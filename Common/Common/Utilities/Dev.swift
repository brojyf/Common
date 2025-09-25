//
//  Dev.swift
//  Common
//
//  Created by 江逸帆 on 9/9/25.
//

#if DEBUG
import Foundation
import SwiftUI

@MainActor
final class PreviewContainer: ObservableObject {
    let session = SessionStore()
    lazy var authVM = AuthVM(session: session)
    
    init (isLoggedIn: Bool = false, page: String = "nil"){
        KCManager.save(.deviceID, "a2332d21-b576-4af5-a4bf-f2fe6f0c9b50")
        if isLoggedIn {
            session.login()
        }
        switch page {
        case "signup":
            authVM.path.append(AuthRoute.verify(email: "email@test.com", scene: .signup))
        case "verify":
            authVM.path.append(AuthRoute.verify(email: "email@test.com", scene: .signup))
            authVM.path.append(AuthRoute.setPassword(email: "email@test.com", scene: .signup))
        default:
            authVM.path = NavigationPath()
        }
    }
}

@MainActor
enum dev {
    static func loggedOut() -> PreviewContainer {
        PreviewContainer()
    }
    static func loggedIn() -> PreviewContainer {
        PreviewContainer(isLoggedIn: true)
    }
    static func verify() -> PreviewContainer {
        PreviewContainer(page: "verify")
    }
}

#endif
