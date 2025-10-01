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
    func setUsernameFirstTime(){
        session.login()
    }
    
    func forgetPasswordWithRouter(){
        path.append(AuthRoute.sendCode(scene: .resetPassword))
    }
    
    func signupWithRouter(){
        path.append(AuthRoute.sendCode(scene: .signup))
    }
    
    func forgetAndResetPassword(){
        withAnimation {
            path = NavigationPath()
        }
    }
    
    func login(email: String, password: String){
        
        guard let deviceID = KCManager.load(.deviceID),
              !deviceID.isEmpty
        else {
            self.errorMsg = "No device id."
            self.hasError = true
            return
        }
        
        svc.login(email: email, pwd: password, deviceID: deviceID)
            .receive(on: DispatchQueue.main)
            .sink(receiveCompletion: { [weak self] completion in
                guard let self else { return }
                NetworkingManager.handleCompletion(completion, &self.hasError, &self.errorMsg)
            }, receiveValue: {  [weak self] resp in
                guard let self else { return }
                print("Store AuthReponse")
                session.login()
                self.resetFlow()
            })
            .store(in: &cancellables)
    }
    
    func createAcctounWithRouter(pwd: String){
        
        guard let deviceID = KCManager.load(.deviceID),
              let token = KCManager.load(.ott),
              !deviceID.isEmpty,
              !token.isEmpty
        else {
            self.errorMsg = "Unknown Error."
            self.hasError = true
            return
        }
        
        svc.createAccount(token: token, password: pwd, deviceID: deviceID)
            .receive(on: DispatchQueue.main)
            .sink(receiveCompletion: { [weak self] completion in
                guard let self else { return }
                NetworkingManager.handleCompletion(completion, &self.hasError, &self.errorMsg)
            }, receiveValue: { [weak self] resp in
                guard let self else { return }
                print("Store AuthReponse")
                path.append(AuthRoute.setUsername)
            })
            .store(in: &cancellables)
    }
    
    func verifyCodeWithRouter(email: String, code: String, scene: AuthScene){
        
        let sceneStr = scene.toString
        
        guard let codeID = KCManager.load(.codeID), !codeID.isEmpty else {
            self.errorMsg = "Please signup again."
            self.hasError = true
            self.resetFlow()
            return
        }
        
        svc.verifyCode(email: email, scene: sceneStr, code: code, codeID: codeID)
            .receive(on: DispatchQueue.main)
            .sink(receiveCompletion: { [weak self] completion in
                guard let self else { return }
                NetworkingManager.handleCompletion(completion, &self.hasError, &self.errorMsg)
            }, receiveValue: { [weak self] resp in
                guard let self else { return }
                KCManager.delete(.codeID)
                KCManager.save(.ott, resp.ott)
                self.path.append(AuthRoute.setPassword(email: email, scene: scene))
            })
            .store(in: &cancellables)
    }
    
    
    func requestCodeWithRouter(email: String, scene: AuthScene, router: Bool = true){
        
        let sceneStr = scene.toString
        
        svc.requestCode(email: email, scene: sceneStr)
            .receive(on: DispatchQueue.main)
            .sink(receiveCompletion: { [weak self] completion in
                guard let self else { return }
                NetworkingManager.handleCompletion(completion, &self.hasError, &self.errorMsg)
            }, receiveValue: { [weak self] resp in
                guard let self else { return }
                KCManager.save(.codeID, "\(resp.codeID)")
                if router {
                    self.path.append(AuthRoute.verify(email: email, scene: scene))
                }
            })
            .store(in: &cancellables)
    }
    
    func resetFlow() { path = .init() }
    
    func dismissError(){
        hasError = false
        errorMsg = nil
    }
    
    func logout(){
        session.logout()
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

