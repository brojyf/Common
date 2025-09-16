//
//  AuthVM.swift
//  Common
//
//  Created by 江逸帆 on 9/9/25.
//

import Foundation
import SwiftUI
import Combine

final class AuthVM: ObservableObject {

    // Route Function = Current View
    @Published var path = NavigationPath()
    
    @Published var hasError: Bool = false
    @Published var errorMsg: String? = nil
    
    private let session: SessionStore
    private let svc: AuthService
    private var cancellables = Set<AnyCancellable>()
    
    init(session: SessionStore) {
        self.session = session
        self.svc = AuthService()
    }
    
    // MARK: - Service Methods
    
    // MARK: - Router Methods
    func login(email: String, password: String){
        session.login()
    }
    
    func setUsernameFirstTime(){
        session.login()
    }
    
    func createAcctounWithRouter(){
        path.append(AuthRoute.setUsername)
    }
    
    func forgetAndResetPassword(){
        withAnimation {
            path = NavigationPath()
        }
    }
    
    func verifyCodeWithRouter(email: String, code: String, scene: AuthScene){
        path.append(AuthRoute.setPassword(email: email, scene: scene))
    }
    
    func requestCodeWithRouter(email: String, scene: AuthScene){
        let sceneStr = scene.toString
        svc.requestCode(email: email, scene: sceneStr)
            .receive(on: DispatchQueue.main)
            .sink(receiveCompletion: { [weak self] completion in
                guard let self else { return }
                NetworkingManager.handleCompletion(completion, &self.hasError, &self.errorMsg)
            }, receiveValue: { [weak self] resp in
                guard let self else { return }
                self.path.append(AuthRoute.verify(email: email, scene: scene))
                print("\(resp.otpID)")
            })
            .store(in: &cancellables)
    }
    
    func forgetPasswordWithRouter(){
        path.append(AuthRoute.sendCode(scene: .resetPassword))
    }
    
    func signupWithRouter(){
        path.append(AuthRoute.sendCode(scene: .signup))
    }
    
    func resetFlow() { path = .init() }
    
    func dismissError(){
        hasError = false
        errorMsg = nil
    }
}

// Mark: -- AuthRoute and AuthScene
enum AuthRoute: Hashable {
    case sendCode(scene: AuthScene)
    case verify(email: String, scene: AuthScene)
    case setPassword(email: String, scene: AuthScene)
    case setUsername
}

enum AuthScene: Hashable {
    case signup, resetPassword
    var toString: String {
        switch self {
        case .signup:
            return "signup"
        case .resetPassword:
            return "reset_password"
        }
    }
}
