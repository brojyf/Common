//
//  API.swift
//  Common
//
//  Created by 江逸帆 on 9/8/25.
//

import Foundation

enum API {
    #if DEBUG
    static let base = URL(string: "http://localhost:8080/api")!
    #else
    static let base = URL(string: "https://patjiang.dpdns.org/api")!
    #endif
    
    enum Auth {
        private static let authBase = base.appendingPathComponent("auth")
                
        static let requestCode    = authBase.appendingPathComponent("request-code")
        static let verifyCode     = authBase.appendingPathComponent("verify-code")
        static let createAccount  = authBase.appendingPathComponent("create-account")
        static let refresh        = authBase.appendingPathComponent("refresh")
        static let forgetPassword = authBase.appendingPathComponent("forget-password")
        static let resetPassword  = authBase.appendingPathComponent("reset-password")
        static let logoutAll      = authBase.appendingPathComponent("logout-all")
    }
    
    enum User {
        private static let userBase = base.appendingPathComponent("user")
        
        static let setUsername  = userBase.appendingPathComponent("set-username")
    }
}
