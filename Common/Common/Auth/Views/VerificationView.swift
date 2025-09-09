//
//  VerificationView.swift
//  Common
//
//  Created by 江逸帆 on 9/9/25.
//

import SwiftUI

struct VerificationView: View {
    
    let email: String
    @State private var code: String = ""
    
    var body: some View {
        VStack {
            Text("Code has been sent to \(email).")
            Text("It'll be expired in 3 minutes.")
            InputField("code", text: $code)
            NavigationLink("Verify"){
                
            }
        }
        .padding()
        .navigationTitle(Text("Verification"))
    }
}

#Preview {
    NavigationStack {
        VerificationView(email: "test@test.com")
    }
}
