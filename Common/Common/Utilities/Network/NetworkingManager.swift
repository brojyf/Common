//
//  NetworkingManager.swift
//  Common
//
//  Created by 江逸帆 on 9/15/25.
//

import Foundation
import Combine

final class NetworkingManager {
    // MARK: - Basic Methods
    static func get(
            url: URL,
            headers: [String: String] = [:],
            timeout: TimeInterval = 30
    ) -> AnyPublisher<Data, NetworkingError>
    {
        var req = URLRequest(url: url)
        headers.forEach { req.addValue($0.value, forHTTPHeaderField: $0.key) }
        
        return URLSession.shared.dataTaskPublisher(for: req)
            .mapError { NetworkingError.transport($0) }
            .subscribe(on: DispatchQueue.global(qos: .userInitiated))
            .tryMap { try handleURLResponse(output: $0, url: url) }
            .mapError { $0 as? NetworkingError ?? .unknown }
            .retryIfTransport(maxRetries: 3, initialDelay: 0.5, jitter: 0.25)
            .eraseToAnyPublisher()
    }
    
    static func post<Request: Encodable>(
        url: URL,
        body: Request,
        headers: [String: String] = [:],
        encoder: JSONEncoder = JSONEncoder(),
        timeout: TimeInterval = 30
    ) -> AnyPublisher<Data, NetworkingError>{
        var req = URLRequest(url: url)
        req.httpMethod = "POST"
        req.addValue("application/json", forHTTPHeaderField: "Content-Type")
        headers.forEach { req.addValue($0.value, forHTTPHeaderField: $0.key) }
        
        do {
            req.httpBody = try encoder.encode(body)
        } catch {
            return Fail(error: NetworkingError.encoding(error)).eraseToAnyPublisher()
        }
        
        return URLSession.shared.dataTaskPublisher(for: req)
            .mapError { NetworkingError.transport($0) }
            .tryMap { try handleURLResponse(output: $0, url: url) }
            .mapError { $0 as? NetworkingError ?? .unknown }
            .retryIfTransport(maxRetries: 3, initialDelay: 0.5, jitter: 0.25)
            .eraseToAnyPublisher()
    }
    
    static func patch<Request: Encodable>(
        url: URL,
        body: Request,
        headers: [String: String] = [:],
        encoder: JSONEncoder = .init(),
        timeout: TimeInterval = 30
    ) -> AnyPublisher<Data, NetworkingError>
    {
        var req = URLRequest(url: url, timeoutInterval: timeout)
        req.httpMethod = "PATCH"
        req.addValue("application/json", forHTTPHeaderField: "Content-Type")
        headers.forEach { req.addValue($0.value, forHTTPHeaderField: $0.key) }

        do {
            req.httpBody = try encoder.encode(body)
        } catch {
            return Fail(error: NetworkingError.encoding(error)).eraseToAnyPublisher()
        }

        return URLSession.shared.dataTaskPublisher(for: req)
            .mapError { NetworkingError.transport($0) }
            .tryMap { try handleURLResponse(output: $0, url: url) }
            .mapError { $0 as? NetworkingError ?? .unknown }
            .retryIfTransport(maxRetries: 3, initialDelay: 0.5, jitter: 0.25)
            .eraseToAnyPublisher()
    }
    
    /// Output -> Data / Log Error
    static func handleURLResponse(
            output: URLSession.DataTaskPublisher.Output,
            url: URL,
            decoder: JSONDecoder = .init()
        ) throws -> Data {
            guard let http = output.response as? HTTPURLResponse else {
                throw NetworkingError.unknown
            }
            if http.statusCode == 204 { return Data() }
            if (200...299).contains(http.statusCode) {
                return output.data
            }

            // 非 2xx：优先尝试解码为 APIErrorBody
            if let api = try? decoder.decode(APIErrorBody.self, from: output.data) {
                logHTTPError(url: url, status: http.statusCode, headers: http.allHeaderFields, data: output.data, tag: "API")
                throw NetworkingError.api(body: api, status: http.statusCode, url: url, headers: http.allHeaderFields, raw: output.data)
            } else {
                // 解码失败，回退为通用 HTTP 错误
                logHTTPError(url: url, status: http.statusCode, headers: http.allHeaderFields, data: output.data, tag: "HTTP")
                throw NetworkingError.http(status: http.statusCode, url: url, headers: http.allHeaderFields, raw: output.data)
            }
    }
    
    static func handleCompletion(_ completion: Subscribers.Completion<NetworkingError>, _ err: inout Bool,_ errMsg: inout String?) {
        switch completion {
        case .finished:
            break
        case .failure(let error):
            print("Error: \(error.localizedDescription)")
            err = true
            errMsg = error.userMessage
        }
    }
    
    // MARK:- Helper Methods
    private static func logHTTPError(
            url: URL,
            status: Int,
            headers: [AnyHashable: Any],
            data: Data,
            tag: String
    ) {
        let headerText = headers
            .map { "\($0.key): \($0.value)" }
            .sorted()
            .joined(separator: "\n")
        
        let preview = data.prefix(8 * 1024)
        let bodyString = String(data: preview, encoding: .utf8)
            ?? preview.map { String(format: "%02x", $0) }.joined(separator: " ")

        print("""
            ⛔️ [NetworkingManager][\(tag)] Error
            ─────────────────────────────────────────────
            URL    : \(url.absoluteString)
            Status : \(status)
            Headers:
            \(headerText)
            Body(<=8KB):
            \(bodyString)
            ─────────────────────────────────────────────
        """)
    }
}


struct APIErrorBody: Decodable, Error {
    let code: String
    let error: String
}

enum NetworkingError: LocalizedError {
    case encoding(Error)
    case api(body: APIErrorBody, status: Int, url: URL, headers: [AnyHashable: Any], raw: Data)
    case http(status: Int, url: URL, headers: [AnyHashable: Any], raw: Data)
    case transport(URLError)
    case unknown

    var userMessage: String {
        switch self {
        case .api(let body, _, _, _, _):
            return body.error
        case .http(let status, _, _, _):
            if status == 401 {
                return "Unauthorized. Please log in again."
            } else if (500...599).contains(status) {
                return "The server is currently unavailable. Please try again later."
            } else {
                return "Request failed (HTTP \(status))."
            }
        case .transport:
            return "Cannot connect to server"
        case .encoding:
            return "Failed to process request. Please try again."
        case .unknown:
            return "An unknown error occurred. Please try again later."
        }
    }

    var debugMessage: String {
        switch self {
        case .api(let body, let status, let url, _, _):
            return "[🔥API] Status: \(status)\nURL: \(url)\nCode: \(body.code)\nError: \(body.error)"
        case .http(let status, let url, _, let raw):
            let preview = String(data: raw.prefix(8*1024), encoding: .utf8)
                ?? "(non-utf8 \(raw.count)B)"
            return "[🔥HTTP] Status: \(status)\nURL: \(url)\nBody(<=8KB):\n\(preview)"
        case .transport(let e):
            return "[📶] Transport: \(e)"
        case .encoding(let error):
            return "[Encoding] Failed to encode data: \(error)"
        case .unknown:
            return "[⚠️] Unknown Error"
        }
    }

    var errorDescription: String? { debugMessage }
}

private extension Publisher where Failure == NetworkingError {
    func retryIfTransport(maxRetries: Int, initialDelay: Double, jitter: Double) -> AnyPublisher<Output, Failure> {
        self.catch { error -> AnyPublisher<Output, Failure> in
            guard case .transport = error, maxRetries > 0 else {
                return Fail(error: error).eraseToAnyPublisher()
            }
            let delay = initialDelay + Double.random(in: -jitter...jitter)
            return self.delay(for: .seconds(delay), scheduler: DispatchQueue.global())
                .retryIfTransport(maxRetries: maxRetries - 1, initialDelay: initialDelay * 2, jitter: jitter)
                .eraseToAnyPublisher()
        }
        .eraseToAnyPublisher()
    }
}
