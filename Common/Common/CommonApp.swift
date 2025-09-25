//
//  CommonApp.swift
//  Common
//
//  Created by 江逸帆 on 9/8/25.
//

import SwiftUI

@main
struct CommonApp: App {
    
    @StateObject private var session: SessionStore
    @StateObject private var authVM: AuthVM
    
    init() {
        KCManager.deviceID()
        let session = SessionStore()
        _session = StateObject(wrappedValue: session)
        _authVM = StateObject(wrappedValue: AuthVM(session: session))
    }
    
    var body: some Scene {
        WindowGroup {
            RootView()
        }
        .environmentObject(session)
        .environmentObject(authVM)
    }
}
