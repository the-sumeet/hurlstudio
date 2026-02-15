# macOS Distribution Guide

## Prerequisites

### 1. Apple Developer Account

You need a paid Apple Developer account ($99/year) at https://developer.apple.com.

### 2. Developer ID Certificate

To distribute outside the Mac App Store you need a **Developer ID Application** certificate.

**Option A: Using Xcode (Recommended)**

1. Open Xcode > Settings > Accounts
2. Sign in with your Apple ID
3. Select your account and click Manage Certificates
4. Click + and choose "Developer ID Application"
5. Xcode creates and installs the certificate and private key automatically

**Option B: Via Apple Developer Portal**

1. Open Keychain Access > Certificate Assistant > Request a Certificate From a Certificate Authority
2. Fill in your email and name, select "Saved to disk", save the `.certSigningRequest` file
3. Go to https://developer.apple.com/account/resources/certificates/list
4. Click + and choose "Developer ID Application"
5. Upload the `.certSigningRequest` file
6. Download the `.cer` file and double-click to install it

**Verify the certificate is installed:**

```
security find-identity -v -p codesigning
```

Expected output:

```
1) ABC123... "Developer ID Application: Your Name (TEAMID)"
   1 valid identities found
```

The certificate must show a private key (key icon in Keychain Access). Without the private key you cannot sign.

### 3. App-Specific Password

Required for notarytool authentication.

1. Go to https://appleid.apple.com > Security > App-Specific Passwords
2. Click + to generate a new one
3. Save it securely — it is shown only once

### 4. create-dmg

```
brew install create-dmg
```

---

## Build

Build a universal binary (Intel + Apple Silicon):

```
wails build -platform darwin/universal
```

Output: `build/bin/hurlstudio.app`

Verify the binary contains both architectures:

```
lipo -info build/bin/hurlstudio.app/Contents/MacOS/hurlstudio
```

Expected:

```
Architectures in the fat file: ... are: x86_64 arm64
```

---

## Sign the App

Get codesigning identifies

```
security find-identity -v -p codesigning

```

Sign with hardened runtime, which is required for notarization:

```
codesign --deep --force --options runtime \
  --sign "Developer ID Application: Your Name (TEAMID)" \
  build/bin/hurlstudio.app
```

Verify:

```
codesign --verify --verbose build/bin/hurlstudio.app
```

Expected:

```
build/bin/hurlstudio.app: valid on disk
build/bin/hurlstudio.app: satisfies its Designated Requirement
```

---

## Notarize the App

### Step 1 - Create a zip

Use `ditto` (not `zip`) to preserve macOS metadata:

```
ditto -c -k --keepParent build/bin/hurlstudio.app hurlstudio_notarize.zip
```

### Step 2 - Submit to Apple

```
xcrun notarytool submit hurlstudio_notarize.zip \
  --apple-id "your@apple.com" \
  --team-id TEAMID \
  --password "xxxx-xxxx-xxxx-xxxx"
```

This returns immediately with a submission ID. Do not use `--wait` — it can appear to hang.

Check status with:

```
xcrun notarytool info <submission-id> \
  --apple-id "your@apple.com" \
  --team-id TEAMID \
  --password "xxxx-xxxx-xxxx-xxxx"
```

If it fails, fetch the log:

```
xcrun notarytool log <submission-id> \
  --apple-id "your@apple.com" \
  --team-id TEAMID \
  --password "xxxx-xxxx-xxxx-xxxx"
```

### Step 3 - Staple

Once status shows `Accepted`:

```
xcrun stapler staple build/bin/hurlstudio.app
```

### Step 4 - Verify

```
spctl --assess --verbose build/bin/hurlstudio.app
```

Expected:

```
build/bin/hurlstudio.app: accepted
source=Notarized Developer ID
```

---

## Create and Notarize a DMG

### Step 1 - Create the DMG

```
create-dmg \
  --volname "Hurl Studio" \
  --volicon "build/appicon.png" \
  --window-pos 200 120 \
  --window-size 600 400 \
  --icon-size 100 \
  --icon "hurlstudio.app" 150 185 \
  --hide-extension "hurlstudio.app" \
  --app-drop-link 450 185 \
  --no-internet-enable \
  "hurlstudio-1.0.0-universal.dmg" \
  "build/bin/"
```

### Step 2 - Sign the DMG

```
codesign --sign "Developer ID Application: Your Name (TEAMID)" \
  hurlstudio-1.0.0-universal.dmg
```

### Step 3 - Notarize the DMG

```
xcrun notarytool submit hurlstudio-1.0.0-universal.dmg \
  --apple-id "your@apple.com" \
  --team-id TEAMID \
  --password "xxxx-xxxx-xxxx-xxxx"
```

Check status, then once `Accepted`:

```
xcrun stapler staple hurlstudio-1.0.0-universal.dmg
```

Verify:

```
spctl --assess --verbose --type open --context context:primary-signature hurlstudio-1.0.0-universal.dmg
```

Expected:

```
hurlstudio-1.0.0-universal.dmg: accepted
source=Notarized Developer ID
```

---

## Storing Credentials in Keychain (Optional)

Instead of passing `--password` every time, store credentials once:

```
xcrun notarytool store-credentials "notarization" \
  --apple-id "your@apple.com" \
  --team-id TEAMID
```

Then use `--keychain-profile "notarization"` in place of `--apple-id`, `--team-id`, and `--password`.

---

## Notes

- If you rebuild the app, you must re-sign, re-notarize, and re-staple everything from scratch.
- The DMG must be notarized separately from the app inside it.
- Stapling embeds the notarization ticket into the file so verification works offline.
- Never use `--wait` with notarytool — submit without it and poll with `notarytool info`.
- Keep app-specific passwords secret and out of scripts. Use keychain profiles instead.
