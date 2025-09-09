//
//  SetUsernameView.swift
//  Common
//
//  Created by 江逸帆 on 9/9/25.
//

import SwiftUI

struct SetUsernameView: View {
    
    @State private var username: String = ""
    
    var body: some View {
        VStack {
            InputField("Username", text: $username)
            Button("Submit"){
                
            }
            
        }
        .padding()
        .navigationBarTitle("Set Username")
    }
}

#Preview {
    NavigationStack {
        SetUsernameView()
    }
}
