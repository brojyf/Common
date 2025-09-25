//
//  AuthService.swift
//  Common
//
//  Created by 江逸帆 on 9/15/25.
//

import Foundation
import Combine

final class AuthService {
    
    func verifyCode(email: String, scene: String, code: String, codeID: String) -> AnyPublisher<OTT, NetworkingError>{
        let body = VerifyCodeBody(email: email, scene: scene, code: code, codeID: codeID)
        return NetworkingManager.post(
            url: API.Auth.verifyCode,
            body: body
        )
        .decode(type: OTT.self, decoder: JSONDecoder())
        .mapError { ($0 as? NetworkingError) ?? .unknown }
        .eraseToAnyPublisher()
    }
    
    func requestCode(email: String, scene: String) -> AnyPublisher<CodeID, NetworkingError>{
        let body = RequestCodeBody(email: email, scene: scene)
        return NetworkingManager.post(
            url: API.Auth.requestCode,
            body: body
        )
        .decode(type: CodeID.self, decoder: JSONDecoder())
        .mapError { ($0 as? NetworkingError) ?? .unknown }
        .eraseToAnyPublisher()
    }
}


struct OTT: Codable {
    let ott: String
    enum CodingKeys: String, CodingKey {
        case ott = "token"
    }
}

struct VerifyCodeBody: Codable {
    let email: String
    let scene: String
    let code: String
    let codeID: String
    enum CodingKeys: String, CodingKey {
        case email, scene, code
        case codeID = "code_id"
    }
}

struct CodeID: Codable {
    let codeID: String
    enum CodingKeys: String, CodingKey {
        case codeID = "code_id"
    }
}

struct RequestCodeBody: Codable {
    let email: String
    let scene: String
}

