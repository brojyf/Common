//
//  SignupView.swift
//  Common
//
//  Created by 江逸帆 on 9/9/25.
//

import SwiftUI

struct SignupView: View {
    
    @State private var email: String = ""
    
    var body: some View {
        VStack {
            InputField("email", text: $email)
            
            NavigationLink("Signup"){
                VerificationView(email: email)
            }
        }
        .padding()
        .navigationTitle(Text("Signup"))
    }
        
}

#Preview {
    NavigationStack {
        SignupView()
    }
}
