//
//  SendCodeView.swift
//  Common
//
//  Created by 江逸帆 on 9/9/25.
//

import SwiftUI

struct SendCodeView: View {
    
    let scene: AuthScene
    @State private var email: String = ""
    @EnvironmentObject var vm: AuthVM
    
    var body: some View {
        VStack {
            InputField("email", text: $email)
            Button("Send Code"){
                vm.requestCodeWithRouter(email: email, scene: scene)
            }
        }
        .padding()
        .navigationTitle(Text(scene == .signup ? "Sign up" : "Reset Password"))
    }
}

#Preview {
    NavigationStack {
        SendCodeView(scene: .signup)
    }
    .environmentObject(AuthVM())
}
#Preview {
    NavigationStack {
        SendCodeView(scene: .resetPassword)
    }
    .environmentObject(AuthVM())
}
