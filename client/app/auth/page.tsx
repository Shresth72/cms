import { googleLogout, useGoogleLogin } from "@react-oauth/google";
import axios from "axios";

type Props = {};

export function isTokenExpired(token: string): [boolean, string] {
  const arrayToken = token.split(".");
  const tokenPayload = JSON.parse(atob(arrayToken[1]));
  console.log(tokenPayload);

  if (Math.floor(new Date().getTime() / 1000) >= tokenPayload?.exp) {
    return [true, tokenPayload?.id as string];
  }
  return [false, ""];
}

export default function UsersPage({}: Props) {
  const login = useGoogleLogin({
    onSuccess: async (codeResponse) => {
      console.log("codeResponse", codeResponse);

      const jwttokenResponse = await axios.get(
        `https://localhost:8000/auth/google/callback?code=${codeResponse.code}`,
      );

      console.log("server signed jwt token: ", jwttokenResponse);
    },
    flow: "auth-code",
  });

  return (
    <>
      <button onClick={() => login()}>Login with Google</button>

      <button onClick={() => googleLogout()}>Logout</button>
    </>
  );
}
