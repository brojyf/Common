//
//  InputField.swift
//  Common
//
//  Created by 江逸帆 on 9/9/25.
//

import SwiftUI

struct InputField: View {
    
    let isSecure: Bool
    let color: Color
    let placeholder: String
    @Binding var text: String
    
    init(isSecure: Bool = false,
         color: Color = .green,
         _ placeholder: String,
         text: Binding<String>) {
        self.isSecure = isSecure
        self.color = color
        self.placeholder = placeholder
        self._text = text
    }
    
    var body: some View {
        Group {
            if isSecure {
                SecureField(placeholder, text: $text)
            } else {
                TextField(placeholder, text: $text)
            }
        }
        .padding()
        .background(isSecure ? .red : color)
    }
}

#Preview {
    InputField("test", text: .constant(""))
}
