"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useAuth } from "@/contexts/AuthContext";
import AuthWrapper from "@/components/auth/AuthWrapper";
import AuthBanner from "@/components/auth/AuthBanner";
import { Rosarivo } from "next/font/google";
import { useState } from "react";
import { LuEye, LuEyeOff } from "react-icons/lu";

const rosarivo = Rosarivo({
  subsets: ["latin"],
  weight: "400",
});

const SignupPage = () => {
  const router = useRouter();
  const { login } = useAuth();

  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [phone, setPhone] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleOnSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setIsLoading(true);

    const API_URL = process.env.NEXT_PUBLIC_API_URL;

    try {
      const payload: {
        name: string;
        email: string;
        password: string;
        phone?: string;
      } = {
        name,
        email,
        password,
      };

      if (phone.trim()) {
        payload.phone = phone;
      }

      const res = await fetch(`${API_URL}/api/v1/auth/register`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });

      const data = await res.json();

      if (!res.ok) {
        if (data.details) {
          const errorMessages = Object.values(data.details).flat().join(", ");
          throw new Error(errorMessages);
        }
        throw new Error(data.message || "Registration failed");
      }

      console.log("Registration successful", data);

      // Extract token and user from response
      const token = data.data?.token || data.token;
      const user = data.data?.user || data.user;

      if (token && user) {
        // Use context to login
        login(token, user);

        // Redirect to home
        router.push("/");
      } else {
        throw new Error("No token or user data received from server");
      }
    } catch (err: unknown) {
      console.error("Registration failed", err);
      setError(
        err instanceof Error
          ? err.message
          : "Something went wrong. Please try again.",
      );
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <AuthWrapper
      left={
        <AuthBanner
          badge="New Member"
          title="Start Your Pilates Journey Today"
          subtitle="Create an account and book your first session with ease."
          image="/images/signup.jpg"
        />
      }
      right={
        <div className="text-background flex flex-col justify-between items-center">
          <h1 className="text-3xl font-bold">DIRO</h1>

          <div className="space-y-10 w-full max-w-sm">
            <div className="text-center">
              <h1 className={`${rosarivo.className} text-3xl`}>
                Create Account
              </h1>
              <p>Sign up to start reserving your pilates sessions.</p>
            </div>

            {error && (
              <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-sm">
                {error}
              </div>
            )}

            <form onSubmit={handleOnSubmit} method="POST" className="space-y-6">
              <div className="flex flex-col gap-1">
                <label className="font-medium">
                  Full Name <span className="text-red-500">*</span>
                </label>
                <input
                  type="text"
                  name="name"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  placeholder="Enter your full name"
                  className="w-full bg-slate-100 p-2 rounded-sm outline-none focus:ring-2 focus:ring-indigo-500"
                  required
                  disabled={isLoading}
                />
              </div>

              <div className="flex flex-col gap-1">
                <label className="font-medium">
                  Email <span className="text-red-500">*</span>
                </label>
                <input
                  type="email"
                  name="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  placeholder="Enter your email"
                  className="w-full bg-slate-100 p-2 rounded-sm outline-none focus:ring-2 focus:ring-indigo-500"
                  required
                  disabled={isLoading}
                />
              </div>

              <div className="flex flex-col gap-1">
                <label className="font-medium">
                  Phone Number{" "}
                  <span className="text-slate-400 text-sm">(Optional)</span>
                </label>
                <input
                  type="tel"
                  name="phone"
                  value={phone}
                  onChange={(e) => setPhone(e.target.value)}
                  placeholder="Enter your phone number"
                  className="w-full bg-slate-100 p-2 rounded-sm outline-none focus:ring-2 focus:ring-indigo-500"
                  disabled={isLoading}
                />
              </div>

              <div className="relative flex flex-col gap-1">
                <label className="font-medium">
                  Password <span className="text-red-500">*</span>
                </label>
                <input
                  type={showPassword ? "text" : "password"}
                  name="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  minLength={6}
                  maxLength={20}
                  required
                  placeholder="Create a password (6-20 characters)"
                  className="w-full bg-slate-100 p-2 rounded-sm outline-none focus:ring-2 focus:ring-indigo-500"
                  disabled={isLoading}
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
              </div>

              <button
                type="submit"
                disabled={isLoading}
                className="w-full bg-indigo-600 text-white py-2 rounded-sm hover:bg-indigo-700 transition disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center"
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
                    Creating account...
                  </>
                ) : (
                  "Sign Up"
                )}
              </button>
            </form>

            <p className="text-xs text-slate-500 text-center">
              By signing up, you agree to our{" "}
              <Link href="/terms" className="underline">
                Terms of Service
              </Link>{" "}
              and{" "}
              <Link href="/privacy" className="underline">
                Privacy Policy
              </Link>
            </p>
          </div>

          <p>
            Already have an account?{" "}
            <Link
              href="/login"
              className="underline text-indigo-600 hover:text-indigo-800"
            >
              Sign In
            </Link>
          </p>
        </div>
      }
    />
  );
};

export default SignupPage;
