//
//  AuthService.swift
//  Common
//
//  Created by 江逸帆 on 9/15/25.
//

import Foundation
import Combine

final class AuthService {
    
    func requestCode(email: String, scene: String) -> AnyPublisher<CodeID, NetworkingError>{
        let body = RequstCodeBody(email: email, scene: scene)
        return NetworkingManager.post(
            url: API.Auth.requestCode,
            body: body
        )
        .decode(type: CodeID.self, decoder: JSONDecoder())
        .mapError { ($0 as? NetworkingError) ?? .unknown }
        .eraseToAnyPublisher()
    }
}

struct CodeID: Codable {
    let codeID: String
    enum CodingKeys: String, CodingKey {
        case codeID = "code_id" // 后端字段是 "code_id"
    }
}

struct RequstCodeBody: Codable {
    let email: String
    let scene: String
}
