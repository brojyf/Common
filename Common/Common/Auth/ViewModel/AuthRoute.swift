//
//  AuthRoute.swift
//  Common
//
//  Created by 江逸帆 on 9/9/25.
//

import Foundation

enum AuthRoute: Hashable {
    case sendCode(scene: AuthScene)
    case verify(email: String, scene: AuthScene)
    case setPassword(email: String, scene: AuthScene)
    case setUsername
}
