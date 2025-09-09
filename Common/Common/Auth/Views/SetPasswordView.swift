//
//  SetPasswordView.swift
//  Common
//
//  Created by 江逸帆 on 9/9/25.
//

import SwiftUI

struct SetPasswordView: View {
    
    @EnvironmentObject var authVM: AuthVM
    
    let scene: AuthScene
    let email: String
    @State private var password: String = ""
    @State private var passwordConfirmation: String = ""
    
    var body: some View {
        VStack {
            InputField(isSecure: true, "password", text: $password)
            InputField(isSecure: true, "Confirm Password", text: $passwordConfirmation)
            
            if scene == .signup {
                Button("Create"){
                    authVM.createAcctounWithRouter()
                }
            } else {
                Button("Reset"){
                    authVM.resetPasswordWihtRouter()
                }
            }
        }
        .padding()
        .navigationTitle(Text(scene == .signup ? "Create Account" : "Reset Password"))
    }
}

#Preview {
    let dev = dev.loggedOut()
    NavigationStack {
        SetPasswordView(scene: .signup, email: "")
    }
    .environmentObject(dev.authVM)
}
