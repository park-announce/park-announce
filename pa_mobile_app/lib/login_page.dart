import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:google_sign_in/google_sign_in.dart';
import 'package:pa_mobile_app/components/pa_login_button.dart';
import 'package:pa_mobile_app/main_page.dart';
import 'package:shared_preferences/shared_preferences.dart';

const List<String> scopes = <String>['email'];

GoogleSignIn _googleSignIn = GoogleSignIn(
  // Optional clientId
  // clientId: 'your-client_id.apps.googleusercontent.com',
  scopes: scopes,
);

class LoginPage extends StatefulWidget {
  const LoginPage({super.key});

  @override
  State<LoginPage> createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {
  @override
  void initState() {
    super.initState();
    _googleSignIn.onCurrentUserChanged.listen((GoogleSignInAccount? account) async {
      // In mobile, being authenticated means being authorized...
      SharedPreferences.getInstance().then((pref) => {pref.setString('Name', account!.displayName!)});
      Navigator.of(context).push(
        MaterialPageRoute<dynamic>(builder: (context) => const MainPage()),
      );
      // However, in the web...
    });

    WidgetsBinding.instance.addPostFrameCallback((timeStamp) {
      showLoginMenu();
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Color(0xFF132555),
      body: Container(
        height: double.infinity,
        width: double.infinity,
        child: Padding(
          padding: const EdgeInsets.all(8),
          child: Stack(children: [
            Positioned(
                bottom: 0,
                right: 0,
                left: 0,
                child: Row(
                  children: [
                    Expanded(
                      child: Padding(
                        padding: const EdgeInsets.all(8),
                        child: Container(
                          decoration: BoxDecoration(borderRadius: BorderRadius.circular(30), color: Colors.white),
                          child: MaterialButton(
                            textColor: Colors.black,
                            child: const Text('Get Started'),
                            onPressed: () {
                              showLoginMenu();
                            },
                          ),
                        ),
                      ),
                    ),
                  ],
                ))
          ]),
        ),
      ),
    );
  }

  Future<dynamic> showLoginMenu() {
    return showModalBottomSheet(
        useSafeArea: true,
        isScrollControlled: true,
        context: context,
        builder: (context) => Container(
              decoration: BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.circular(30),
              ),
              padding: const EdgeInsets.only(bottom: 20, top: 10),
              width: double.infinity,
              child: Padding(
                padding: const EdgeInsets.all(18),
                child: Wrap(
                  children: [
                    PaLoginButton(
                        backColor: Colors.blue,
                        textColor: Colors.white,
                        onPressedFunction: () {},
                        child: const Text(
                          'Sign Up',
                          style: TextStyle(fontSize: 11),
                        )),
                    PaLoginButton(
                        backColor: Colors.white,
                        textColor: Colors.black,
                        onPressedFunction: () {},
                        child: const Text(
                          'Log In',
                          style: TextStyle(fontSize: 11),
                        )),
                    const Row(
                      children: [
                        Expanded(
                          child: Divider(
                            color: Colors.grey,
                            height: 10,
                          ),
                        ),
                        Text('Or', style: TextStyle(fontSize: 10)),
                        Expanded(
                          child: Divider(
                            color: Colors.grey,
                            height: 10,
                          ),
                        ),
                      ],
                    ),
                    PaLoginButton(
                      backColor: Colors.white,
                      textColor: Colors.black,
                      onPressedFunction: () {},
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          const FaIcon(
                            FontAwesomeIcons.apple,
                            size: 18,
                          ),
                          const Text(
                            'Continue With Apple',
                            style: TextStyle(fontSize: 11),
                          ),
                          Container()
                        ],
                      ),
                    ),
                    PaLoginButton(
                      backColor: Colors.white,
                      textColor: Colors.black,
                      onPressedFunction: () {},
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          const FaIcon(
                            FontAwesomeIcons.facebook,
                            color: Colors.blue,
                            size: 18,
                          ),
                          const Text(
                            'Continue With Facebook',
                            style: TextStyle(fontSize: 11),
                          ),
                          Container()
                        ],
                      ),
                    ),
                    PaLoginButton(
                      backColor: Colors.white,
                      textColor: Colors.black,
                      onPressedFunction: () {
                        _googleSignIn.signIn();
                      },
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          const FaIcon(FontAwesomeIcons.google, size: 18),
                          const Text(
                            'Continue With Google',
                            style: TextStyle(fontSize: 11),
                          ),
                          Container()
                        ],
                      ),
                    )
                  ],
                ),
              ),
            ));
  }
}
