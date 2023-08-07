import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:google_sign_in/google_sign_in.dart';
import 'package:jwt_decoder/jwt_decoder.dart';
import 'package:pa_mobile_app/components/pa_button.dart';
import 'package:pa_mobile_app/components/pa_login_button.dart';
import 'package:pa_mobile_app/models/api_error_response.dart';
import 'package:pa_mobile_app/models/check_api_token_response.dart';
import 'package:pa_mobile_app/models/preregister_google_response.dart';
import 'package:pa_mobile_app/pages/login_page.dart';
import 'package:pa_mobile_app/pages/map_page.dart';
import 'package:pa_mobile_app/pages/register_mail_page.dart';
import 'package:pa_mobile_app/pages/register_page.dart';
import 'package:pa_mobile_app/service.dart';
import 'package:pa_mobile_app/utils/navigation_utils.dart' as nav_utils;
import 'package:shared_preferences/shared_preferences.dart';

const List<String> scopes = <String>['email'];

GoogleSignIn _googleSignIn = GoogleSignIn(
  // Optional clientId
  // clientId: 'your-client_id.apps.googleusercontent.com',
  scopes: scopes,
);

class OnBoardingPage extends StatefulWidget {
  const OnBoardingPage({super.key});

  @override
  State<OnBoardingPage> createState() => _OnBoardingPageState();
}

class _OnBoardingPageState extends State<OnBoardingPage> {
  @override
  void initState() {
    super.initState();
    //_googleSignIn.signInSilently();
    _googleSignIn.onCurrentUserChanged.listen((GoogleSignInAccount? account) async {
      _handleGoogleAccount(account!);
    });

    WidgetsBinding.instance.addPostFrameCallback((timeStamp) {
      //showLoginMenu();
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Padding(
        padding: EdgeInsets.all(20),
        child: SizedBox(
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
                        child: PaButton(
                          text: 'Get Started',
                          onPressedFunction: () {
                            showLoginMenu();
                          },
                        ),
                      ),
                    ],
                  ))
            ]),
          ),
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
              width: double.infinity,
              child: Padding(
                padding: const EdgeInsets.all(18),
                child: Wrap(
                  children: [
                    PaLoginButton(
                        onPressedFunction: () {
                          nav_utils.navigate(context, const RegisterMailPage());
                        },
                        text: 'Sign Up'),
                    PaLoginButton(
                      onPressedFunction: () {
                        nav_utils.navigate(context, const LoginPage());
                      },
                      text: 'Log In',
                    ),
                    Row(
                      children: [
                        const Expanded(
                          child: Divider(
                            color: Colors.grey,
                            height: 10,
                          ),
                        ),
                        Text('Or', style: Theme.of(context).textTheme.bodySmall),
                        const Expanded(
                          child: Divider(
                            color: Colors.grey,
                            height: 10,
                          ),
                        ),
                      ],
                    ),
                    PaLoginButton(
                      onPressedFunction: () {
                        _googleSignIn.signIn().then((value) {
                          _handleGoogleAccount(value!);
                        });
                      },
                      text: 'Continue With Google',
                      icon: FaIcon(FontAwesomeIcons.google),
                    )
                  ],
                ),
              ),
            ));
  }

  void _handleGoogleAccount(GoogleSignInAccount account) {
    {
      // In mobile, being authenticated means being authorized...
      SharedPreferences.getInstance().then((pref) {
        pref.setString('Email', 'undefined');
        account!.authentication.then((value) {
          checkApiToken(value.idToken!).then((apiTokenResponseBody) {
            if (apiTokenResponseBody is CheckApiTokenResponse) {
              final decoded = JwtDecoder.decode(apiTokenResponseBody.token);
              pref.setString('Email', decoded["email"].toString());
              pref.setString('Token', apiTokenResponseBody.token);
              pref.setString('IdToken', value.idToken!);
              pref.setString('Name', account.displayName!);
              nav_utils.navigate(context, const MapPage());
            } else if (apiTokenResponseBody is ApiErrorResponse && apiTokenResponseBody.code == 'exp.user.notfound') {
              registerGoogle(value.idToken!).then((apiTokenResponseBody) {
                if (apiTokenResponseBody is PreRegisterGoogleResponse) {
                  /*
                  final decoded = JwtDecoder.decode(apiTokenResponseBody.token);
                  pref.setString('Email', decoded["email"].toString());
                  pref.setString('Token', apiTokenResponseBody.token);
                  pref.setString('IdToken', value.idToken!);
                  pref.setString('Name', account.displayName!);
                  nav_utils.navigate(context, const MapPage());
                  */
                  nav_utils.navigate(
                      context, RegisterPage(email: account.email, firstName: account.displayName!.split(' ')[0], lastName: account.displayName!.split(' ')[1]));
                } else {
                  //TODO:Hata ver
                }
              });
            }
          });
        });
      });

      // However, in the web...
    }
  }
}
