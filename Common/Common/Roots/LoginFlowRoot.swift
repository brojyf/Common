//
//  LoginFlowRoot.swift
//  Common
//
//  Created by 江逸帆 on 9/9/25.
//

import SwiftUI

struct LoginFlowRoot: View {
    @EnvironmentObject var authVM: AuthVM
    
    var body: some View {
        NavigationStack(path: $authVM.path){
            LoginView()
                .navigationDestination(for: AuthRoute.self){ route in
                    switch route {
                    case .sendCode(let scene):
                        SendCodeView(scene: scene)
                    case .verify(let email, let scene):
                        VerificationView(email: email, scene: scene)
                    case .setPassword(email: let email, scene: let scene):
                        SetPasswordView(scene: scene, email: email)
                    case .setUsername:
                        SetUsernameView()
                    }
                }
        }
        .task { authVM.resetFlow() }
    }
}

#Preview {
    let dev = dev.loggedOut()
    LoginFlowRoot()
        .environmentObject(dev.authVM)
}
