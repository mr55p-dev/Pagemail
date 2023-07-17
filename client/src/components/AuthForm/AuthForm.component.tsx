import { DataState } from "../../lib/data";
import { pb, useUser } from "../../lib/pocketbase";
// import signinUrl from "../../assets/google-auth/2x/btn_google_signin_light_normal_web@2x.png";
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
import { UserRecord } from "../../lib/datamodels";
import React from "react";
import { NotificationCtx } from "../../lib/notif";

interface LoginForm {
  email: string;
  password: string;
}

export const Login = () => {
  const [reqState, setReqState] = React.useState<DataState>(DataState.UNKNOWN);
  const { login, authErr } = useUser();
  const { notifErr } = React.useContext(NotificationCtx);
  const { register, handleSubmit } = useForm<LoginForm>();

  function onSubmit(data: LoginForm) {
    setReqState(DataState.PENDING);
    login(async () => {
      const response = await pb
        .collection("users")
        .authWithPassword<UserRecord>(data.email, data.password);
      return response.record;
    })
      .then(() => setReqState(DataState.SUCCESS))
      .catch((e) => {
        setReqState(DataState.FAILED);
        notifErr("Error logging in", e);
      });
  }

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
            <Button type="submit" disabled={reqState === DataState.PENDING}>
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
  const { login } = useUser();
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
      const res = await pb
        .collection("users")
        .authWithPassword<UserRecord>(data.email, data.password);
      return res.record;
    });

  return (
    <>
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

// const GoogleAuth = () => {
//   const { login } = useUser();
//   const handleGoogle = () => {
//     login(async () => {
//       await pb.collection("users").authWithOAuth2({ provider: "google" });
//     });
//   };
//
//   return (
//     <Button sx={{ mx: 2 }} onClick={handleGoogle}>
//       <img src={signinUrl} width="200px" />
//     </Button>
//   );
// };
