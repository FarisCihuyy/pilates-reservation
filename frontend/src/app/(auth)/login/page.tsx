"use client";

import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { useAuth } from "@/contexts/AuthContext";
import { useBooking } from "@/contexts/BookingContext";
import AuthWrapper from "@/components/auth/AuthWrapper";
import AuthBanner from "@/components/auth/AuthBanner";
import { Rosarivo } from "next/font/google";
import { LuEye, LuEyeOff } from "react-icons/lu";
import { useState, useEffect } from "react";

const rosarivo = Rosarivo({ subsets: ["latin"], weight: "400" });

const LoginPage = () => {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { login } = useAuth();
  const { loadFromSessionStorage } = useBooking();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [rememberMe, setRememberMe] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);

  // Check for success messages from query params
  useEffect(() => {
    const registered = searchParams.get("registered");
    if (registered === "true") {
      setSuccessMessage("Registration successful! Please sign in.");
    }
  }, [searchParams]);

  const handleOnSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setIsLoading(true);

    const API_URL = process.env.NEXT_PUBLIC_API_URL;

    try {
      const res = await fetch(`${API_URL}/api/v1/auth/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password }),
      });

      const data = await res.json();
      console.log(data);

      if (!res.ok) {
        throw new Error(data.details || "Login failed");
      }

      console.log("Login successful", data);

      // Extract token and user from response
      const token = data.data?.token || data.token;
      const user = data.data?.user || data.user;

      if (token && user) {
        // Use context to login
        login(token, user);

        // Restore booking data if exists
        loadFromSessionStorage();

        // Get redirect URL from query params
        const redirectUrl = searchParams.get("redirect") || "/";

        // Redirect to intended page or home
        router.push(redirectUrl);
      } else {
        throw new Error("No token or user data received from server");
      }
    } catch (err: unknown) {
      console.error("Login failed", err);
      setError(
        err instanceof Error
          ? err.message
          : "Invalid email or password. Please try again.",
      );
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <AuthWrapper
      left={
        <AuthBanner
          badge="Member Access"
          title="Manage Your Bookings Effortlessly"
          subtitle="Access your schedule, manage reservations, and stay in control — all in one place."
          image="/images/login.jpg"
        />
      }
      right={
        <div className="text-background flex flex-col justify-between items-center">
          <h1 className="text-3xl font-bold">DIRO</h1>

          <div className="space-y-10 w-full max-w-sm">
            <div className="text-center">
              <h1 className={`${rosarivo.className} text-3xl`}>Welcome Back</h1>
              <p>Sign in to manage your reservations and bookings.</p>
            </div>

            {/* Success Message */}
            {successMessage && (
              <div className="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded-sm">
                {successMessage}
              </div>
            )}

            {/* Error Message */}
            {error && (
              <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-sm">
                {error}
              </div>
            )}

            <form onSubmit={handleOnSubmit} className="space-y-6">
              <div className="flex flex-col gap-1">
                <label htmlFor="email" className="font-medium">
                  Email <span className="text-red-500">*</span>
                </label>
                <input
                  type="email"
                  id="email"
                  name="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  placeholder="Enter your email"
                  required
                  disabled={isLoading}
                  className="bg-slate-100 p-2 rounded-sm outline-none focus:ring-2 focus:ring-indigo-500 disabled:opacity-50"
                />
              </div>

              <div className="relative flex flex-col gap-1">
                <label htmlFor="password" className="font-medium">
                  Password <span className="text-red-500">*</span>
                </label>
                <input
                  type={showPassword ? "text" : "password"}
                  id="password"
                  name="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  placeholder="Enter your password"
                  required
                  disabled={isLoading}
                  className="bg-slate-100 p-2 rounded-sm outline-none focus:ring-2 focus:ring-indigo-500 disabled:opacity-50"
                />
                <button
                  type="button"
                  className="absolute right-2 top-10"
                  onClick={() => setShowPassword(!showPassword)}
                  disabled={isLoading}
                >
                  {showPassword ? (
                    <LuEye className="text-xl opacity-70 cursor-pointer hover:opacity-100" />
                  ) : (
                    <LuEyeOff className="text-xl opacity-70 cursor-pointer hover:opacity-100" />
                  )}
                </button>

                <div className="flex justify-between items-center mt-2">
                  <div className="flex items-center gap-2">
                    <input
                      type="checkbox"
                      id="remember"
                      checked={rememberMe}
                      onChange={(e) => setRememberMe(e.target.checked)}
                      disabled={isLoading}
                      className="cursor-pointer"
                    />
                    <label
                      htmlFor="remember"
                      className="text-sm cursor-pointer"
                    >
                      Remember Me
                    </label>
                  </div>
                  <Link
                    href="/forgot-password"
                    className="text-sm underline text-indigo-600 hover:text-indigo-800"
                  >
                    Forgot Password?
                  </Link>
                </div>
              </div>

              <button
                type="submit"
                disabled={isLoading}
                className="w-full bg-indigo-600 text-white py-2 rounded-sm hover:bg-indigo-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center"
              >
                {isLoading ? (
                  <>
                    <svg
                      className="animate-spin -ml-1 mr-3 h-5 w-5 text-white"
                      xmlns="http://www.w3.org/2000/svg"
                      fill="none"
                      viewBox="0 0 24 24"
                    >
                      <circle
                        className="opacity-25"
                        cx="12"
                        cy="12"
                        r="10"
                        stroke="currentColor"
                        strokeWidth="4"
                      ></circle>
                      <path
                        className="opacity-75"
                        fill="currentColor"
                        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                      ></path>
                    </svg>
                    Signing in...
                  </>
                ) : (
                  "Sign In"
                )}
              </button>
            </form>

            {/* Divider */}
            <div className="relative">
              <div className="absolute inset-0 flex items-center">
                <div className="w-full border-t border-slate-300"></div>
              </div>
              <div className="relative flex justify-center text-sm">
                <span className="px-2 bg-white text-slate-500">
                  Or continue with
                </span>
              </div>
            </div>

            {/* Social Login (Optional) */}
            <div className="grid grid-cols-2 gap-4">
              <button
                type="button"
                className="flex items-center justify-center gap-2 border border-slate-300 rounded-sm py-2 hover:bg-slate-50 transition-colors"
                disabled={isLoading}
              >
                <svg className="w-5 h-5" viewBox="0 0 24 24">
                  <path
                    fill="currentColor"
                    d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
                  />
                  <path
                    fill="currentColor"
                    d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
                  />
                  <path
                    fill="currentColor"
                    d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
                  />
                  <path
                    fill="currentColor"
                    d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
                  />
                </svg>
                <span className="text-sm">Google</span>
              </button>

              <button
                type="button"
                className="flex items-center justify-center gap-2 border border-slate-300 rounded-sm py-2 hover:bg-slate-50 transition-colors"
                disabled={isLoading}
              >
                <svg
                  className="w-5 h-5"
                  fill="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path d="M24 12.073c0-6.627-5.373-12-12-12s-12 5.373-12 12c0 5.99 4.388 10.954 10.125 11.854v-8.385H7.078v-3.47h3.047V9.43c0-3.007 1.792-4.669 4.533-4.669 1.312 0 2.686.235 2.686.235v2.953H15.83c-1.491 0-1.956.925-1.956 1.874v2.25h3.328l-.532 3.47h-2.796v8.385C19.612 23.027 24 18.062 24 12.073z" />
                </svg>
                <span className="text-sm">Facebook</span>
              </button>
            </div>
          </div>

          <p>
            Don&apos;t have an account?{" "}
            <Link
              href="/signup"
              className="underline text-indigo-600 hover:text-indigo-800"
            >
              Sign Up
            </Link>
          </p>
        </div>
      }
    />
  );
};

export default LoginPage;
