//
//  CeateAccountView.swift
//  Common
//
//  Created by 江逸帆 on 9/9/25.
//

import SwiftUI

struct CeateAccountView: View {
    
    @State private var password: String = ""
    @State private var passwordConfirmation: String = ""
    
    var body: some View {
        VStack {
            InputField(isSecure: true, "password", text: $password)
            InputField(isSecure: true, "Confirm Password", text: $passwordConfirmation)
            
            NavigationLink("Create") {
                SetUsernameView()
            }
        }
        .padding()
        .navigationTitle(Text("Create Account"))
    }
}

#Preview {
    NavigationStack {
        CeateAccountView()
    }
}
