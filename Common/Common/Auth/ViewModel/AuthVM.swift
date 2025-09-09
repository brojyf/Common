//
//  AuthVM.swift
//  Common
//
//  Created by 江逸帆 on 9/9/25.
//

import Foundation
import SwiftUI

final class AuthVM: ObservableObject {

    @Published var path = NavigationPath()
    
    // MARK: - Router Methods
    func login(email: String, password: String){}
    
    func createAcctounWithRouter(){
        path.append(AuthRoute.setUsername)
    }
    func resetPasswordWihtRouter(){
        withAnimation {
            path = NavigationPath()
        }
    }
    
    func verifyCodeWithRouter(email: String, code: String, scene: AuthScene){
        path.append(AuthRoute.setPassword(email: email, scene: scene))
    }
    
    func requestCodeWithRouter(email: String, scene: AuthScene){
        path.append(AuthRoute.verify(email: email, scene: scene))
    }
    
    func forgotPasswordWithRouter(){
        path.append(AuthRoute.sendCode(scene: .resetPassword))
    }
    
    func signupWithRouter(){
        path.append(AuthRoute.sendCode(scene: .signup))
    }
}
