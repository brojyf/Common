//
//  SessionStore.swift
//  Common
//
//  Created by 江逸帆 on 9/9/25.
//

import Foundation

final class SessionStore: ObservableObject {
    @Published private(set) var isLoggedIn: Bool = false
    @Published private(set) var userID: UInt64?
    @Published private(set) var accessToken: String?
    @Published private(set) var refreshToken: String?
    
    func login(){
        self.isLoggedIn = true
    }
    
    func logout(){
        clear()
        self.isLoggedIn = false
    }
    
    
    func clear(){
        isLoggedIn = false
        userID = nil
        accessToken = nil
        refreshToken = nil
        clearKeychain()
    }
    
    private func clearKeychain(){
        
    }
}
