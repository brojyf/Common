//
//  VerificationView.swift
//  Common
//
//  Created by 江逸帆 on 9/9/25.
//

import SwiftUI

struct VerificationView: View {
    
    @EnvironmentObject var authVM: AuthVM
    
    let email: String
    let scene: AuthScene
    @State private var code: String = ""

    var body: some View {
        VStack {
            Text("Code has been sent to \(email).")
            Text("It'll be expired in 3 minutes.")
            InputField("code", text: $code)
            Button("Verify"){
                authVM.verifyCodeWithRouter(email: email, code: code, scene: scene)
            }
        }
        .padding()
        .navigationTitle(Text("Verification"))
    }
}

#Preview {
    let dev = dev.loggedIn()
    NavigationStack {
        VerificationView(email: "test@test.com", scene: .signup)
    }
    .environmentObject(dev.authVM)
}
