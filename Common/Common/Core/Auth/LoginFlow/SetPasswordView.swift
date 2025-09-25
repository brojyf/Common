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
            
            smartButton
        }
        .padding()
        .navigationTitle(Text(scene == .signup ? "Create Account" : "Reset Password"))
        .alert(isPresented: $authVM.hasError){
            Alert(
                title: Text("Error"),
                message: Text(authVM.errorMsg ?? "Unknown Error"),
                dismissButton: .default(Text("OK")){
                    authVM.dismissError()
                }
            )
        }
    }
}

#Preview {
    let dev = dev.loggedOut()
    NavigationStack {
        SetPasswordView(scene: .signup, email: "")
    }
    .environmentObject(dev.authVM)
}

// MARK: - Extension
extension SetPasswordView {
    private var smartButton: some View {
        if scene == .signup {
            Button("Create"){
                if isSamePwd() {
                    authVM.createAcctounWithRouter(pwd: password)
                } else {
                    authVM.errorMsg = "Password not match"
                    authVM.hasError = true
                }
            }
        } else {
            Button("Reset"){
                if isSamePwd() {
                    authVM.forgetAndResetPassword()
                } else {
                    
                }
            }
        }
    }
    
    private func isSamePwd() -> Bool {
        return password == passwordConfirmation
    }
}
