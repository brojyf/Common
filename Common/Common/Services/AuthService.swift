//
//  AuthService.swift
//  Common
//
//  Created by 江逸帆 on 9/15/25.
//

import Foundation
import Combine

final class AuthService {
    
    func requestCode(email: String, scene: String) -> AnyPublisher<RequestCodeResponse, NetworkingError>{
        let body = RequstCodeBody(email: email, scene: scene)
        return NetworkingManager.post(
            url: API.Auth.requestCode,
            body: body
        )
        .decode(type: RequestCodeResponse.self, decoder: JSONDecoder())
        .mapError { ($0 as? NetworkingError) ?? .unknown }
        .eraseToAnyPublisher()
    }
}

struct RequestCodeResponse: Codable {
    let otpID: String
    enum CodingKeys: String, CodingKey {
        case otpID = "otp_id"
    }
}

struct RequstCodeBody: Codable {
    let email: String
    let scene: String
}
