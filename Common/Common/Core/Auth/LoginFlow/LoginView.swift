//
//  LoginView.swift
//  Common
//
//  Created by 江逸帆 on 9/8/25.
//

import SwiftUI

struct LoginView: View {
    
    @EnvironmentObject var vm: AuthVM
    
    @State private var email: String = ""
    @State private var password: String = ""
    
    var body: some View {
        VStack {
            Text("Welcome to Common")
            inputSection
            buttonsSection
        }
        .padding()
        .alert(isPresented: $vm.hasError){
            Alert(
                title: Text("Error"),
                message: Text(vm.errorMsg ?? "Unknown Error"),
                dismissButton: .default(Text("OK")){
                    vm.dismissError()
                }
            )
        }
    }
}

// MARK: - Preview
#Preview {
    let dev = dev.loggedIn()
    NavigationStack {
        LoginView()
    }
    .environmentObject(dev.authVM)
}

// MARK: - Extension
extension LoginView {
    private var buttonsSection: some View {
        VStack {
            Button("Login") {
                withAnimation(.spring){
                    vm.login(email: email, password: password)
                }
            }
            
            HStack {
                Button("Forget Password"){ vm.forgetPasswordWithRouter() }
                Spacer()
                Button("Signup"){ vm.signupWithRouter() }
            }
        }
    }
    
    private var inputSection: some View {
        VStack {
            InputField("email", text: $email)
            InputField(isSecure: true, "password", text: $password)
        }
    }
}
