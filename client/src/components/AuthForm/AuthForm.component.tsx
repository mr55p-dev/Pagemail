// import React, { ChangeEventHandler } from "react";
import { AuthState } from "../../lib/data";
import { pb, useUser } from "../../lib/pocketbase";
import signinUrl from "../../assets/google-auth/2x/btn_google_signin_light_normal_web@2x.png";
// import { useNotification } from "../../lib/notif";
import { useForm } from "react-hook-form";
import {
  Typography,
  Link,
  Button,
  FormControl,
  FormLabel,
  Input,
  Checkbox,
  FormHelperText,
  Stack,
} from "@mui/joy";

interface LoginForm {
  email: string;
  password: string;
}

export const Login = () => {
  const { login, authState, authErr } = useUser();
  const { register, handleSubmit } = useForm<LoginForm>();

  const onSubmit = (data: LoginForm) =>
    login(() =>
      pb.collection("users").authWithPassword(data.email, data.password)
    );

  const handlePasswordReset = () => {
    alert(
      "Password reset is not yet available, please contact ellis@pagemail.io for assistance"
    );
  };

  return (
    <>
      {authErr ? <div>{authErr.message}</div> : undefined}
      <form onSubmit={handleSubmit(onSubmit)}>
        <Stack spacing={2}>
          <FormControl>
            <FormLabel>Email</FormLabel>
            <Input type="email" {...register("email", { required: true })} />
          </FormControl>
          <FormControl>
            <FormLabel>Password</FormLabel>
            <Input type="password" {...register("password")} />
          </FormControl>
          <FormControl>
            <Button type="submit" disabled={authState === AuthState.PENDING}>
              Sign in
            </Button>
          </FormControl>
          <FormControl>
            <Typography fontSize="sm" sx={{ alignSelf: "center" }}>
              <Link onClick={handlePasswordReset}>Reset password</Link>
            </Typography>
          </FormControl>
        </Stack>
      </form>
    </>
  );
};

interface SignupForm {
  email: string;
  password: string;
  passwordConfirm: string;
  name?: string;
  subscribed: boolean;
}

export const SignUp = () => {
  const { login, authErr } = useUser();
  // const { trigger, component } = useNotification();
  const {
    register,
    handleSubmit,
    formState: { errors },
    watch,
  } = useForm<SignupForm>();

  const onSubmit = (data: SignupForm) =>
    login(async () => {
      await pb.collection("users").create(data);
      await pb.collection("users").authWithPassword(data.email, data.password);
    });

  return (
    <>
      {authErr ? <div>{authErr.message}</div> : undefined}
      <form onSubmit={handleSubmit(onSubmit, (data) => console.error(data))}>
        <Stack spacing={1}>
          <FormControl>
            <FormLabel>Name</FormLabel>
            <Input type="text" {...register("name")} />
          </FormControl>
          <FormControl>
            <FormLabel>Email</FormLabel>
            <Input type="email" {...register("email", { required: true })} />
          </FormControl>
          <FormControl>
            <FormLabel>Password</FormLabel>
            <Input
              type="password"
              color={
                errors.password || errors.passwordConfirm ? "danger" : "neutral"
              }
              {...register("password", { required: true })}
            />
          </FormControl>
          <FormControl>
            <FormLabel>Repeat password</FormLabel>
            <Input
              type="password"
              color={errors.passwordConfirm ? "danger" : "neutral"}
              {...register("passwordConfirm", {
                required: true,
                validate: (val: string) => {
                  if (watch("password") != val) {
                    return "Passwords do not match";
                  }
                },
              })}
            />
            {errors.passwordConfirm && (
              <FormHelperText color="danger">
                {errors.passwordConfirm.message}
              </FormHelperText>
            )}
          </FormControl>
          <FormControl>
            <Checkbox
              defaultChecked
              label="Subscribe?"
              {...register("subscribed")}
            />
            <FormHelperText>
              You'll receive a briefing each morning of yesterdays pages
            </FormHelperText>
          </FormControl>
          <FormControl>
            <Button type="submit">Sign Up</Button>
          </FormControl>
        </Stack>
      </form>
    </>
  );
};

const GoogleAuth = () => {
  const { login } = useUser();
  const handleGoogle = () => {
    login(async () => {
      await pb.collection("users").authWithOAuth2({ provider: "google" });
    });
  };

  return (
    <Button sx={{ mx: 2 }} onClick={handleGoogle}>
      <img src={signinUrl} width="200px" />
    </Button>
  );
};
