//
//  LoginView.swift
//  Common
//
//  Created by 江逸帆 on 9/8/25.
//

import SwiftUI

struct LoginView: View {
    
    @State private var email: String = ""
    @State private var password: String = ""
    @EnvironmentObject var vm: AuthVM
    
    var body: some View {
        VStack {
            Text("Welcome to Common")
            inputSection
            buttonsSection
        }
        .padding()
    }
}

// MARK: - Preview
#Preview {
    NavigationStack {
        LoginView()
    }
    .environmentObject(AuthVM())
}

// MARK: - Extension
extension LoginView {
    private var buttonsSection: some View {
        VStack {
            Button("Login") {
                vm.login(email: email, password: password)
            }
            
            HStack {
                Button("Forget Password"){ vm.forgotPasswordWithRouter() }
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
