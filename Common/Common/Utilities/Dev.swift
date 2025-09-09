//
//  Dev.swift
//  Common
//
//  Created by 江逸帆 on 9/9/25.
//

#if DEBUG
import Foundation

@MainActor
final class PreviewContainer: ObservableObject {
    let session = SessionStore()
    lazy var authVM = AuthVM(session: session)
    
    init (isLoggedIn: Bool = false){
        if isLoggedIn {
            session.login()
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
}

#endif
