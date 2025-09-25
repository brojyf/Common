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
            Text("A Code has been sent to \(email).")
            Text("It'll be expired in 3 minutes.")
            Text("Only the latest code is valid")
            HStack {
                InputField("code", text: $code)
                Button("Send"){
                    authVM.requestCodeWithRouter(email: email, scene: scene, router: false)
                }
            }
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
