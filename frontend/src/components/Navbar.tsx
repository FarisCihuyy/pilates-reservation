"use client";

import clsx from "clsx";
import Link from "next/link";
import { RxHamburgerMenu } from "react-icons/rx";
import { IoClose } from "react-icons/io5";
import { FiUser, FiLogOut } from "react-icons/fi";
import { useState } from "react";
import { useAuth } from "@/contexts/AuthContext";
import { useRouter } from "next/navigation";

const Navbar = ({ className }: { className?: string }) => {
  const { user, isAuthenticated, logout, isLoading } = useAuth();
  const router = useRouter();
  const [isOpen, setIsOpen] = useState(false);
  const [showUserMenu, setShowUserMenu] = useState(false);

  const handleLogout = () => {
    logout();
    setShowUserMenu(false);
    setIsOpen(false);
    router.push("/");
  };

  const handleLogin = () => {
    setIsOpen(false);
    router.push("/login");
  };

  return (
    <nav
      className={`fixed z-40 w-full flex justify-between items-center p-4 md:p-8 uppercase font-medium ${className}`}
    >
      <Link href="/" className="font-bold text-2xl tracking-widest md:hidden">
        Diro
      </Link>
      <RxHamburgerMenu
        onClick={() => setIsOpen(true)}
        className="md:hidden text-2xl cursor-pointer"
      />
      <div
        className={clsx(
          "fixed inset-0 bg-background flex flex-col md:flex-row justify-center md:justify-between w-full gap-12 items-center text-2xl md:text-base lg:text-lg md:translate-x-0 md:bg-transparent md:static transition-transform",
          { "translate-x-0": isOpen, "translate-x-full": !isOpen },
        )}
      >
        <IoClose
          onClick={() => setIsOpen(false)}
          className="text-3xl absolute top-4 right-4 md:hidden cursor-pointer"
        />
        <div className="flex flex-col items-center md:flex-row gap-12">
          <Link href="/schedule" onClick={() => setIsOpen(false)}>
            Schedule
          </Link>
          <Link href="/studios" onClick={() => setIsOpen(false)}>
            Studios
          </Link>
        </div>
        <Link
          href="/"
          onClick={() => setIsOpen(false)}
          className="hidden md:inline-block absolute left-1/2 -translate-x-1/2 font-bold text-3xl tracking-widest"
        >
          Diro
        </Link>
        <div className="flex flex-col items-center md:flex-row gap-12">
          <Link href="#pricing" onClick={() => setIsOpen(false)}>
            Pricing
          </Link>

          {/* Loading State */}
          {isLoading ? (
            <div className="w-8 h-8 border-2 border-current border-t-transparent rounded-full animate-spin md:w-6 md:h-6"></div>
          ) : isAuthenticated && user ? (
            /* Authenticated User */
            <div className="relative">
              {/* Desktop User Menu */}
              <button
                onClick={() => setShowUserMenu(!showUserMenu)}
                className="hidden md:flex items-center gap-2 text-foreground font-bold hover:opacity-80 transition-opacity"
              >
                <FiUser className="text-xl" />
                <span className="normal-case">{user.name}</span>
              </button>

              {/* Desktop Dropdown */}
              {showUserMenu && (
                <>
                  <div
                    className="fixed inset-0 z-10"
                    onClick={() => setShowUserMenu(false)}
                  ></div>
                  <div className="absolute right-0 top-full mt-2 w-48 bg-white dark:bg-slate-800 rounded-lg shadow-lg py-2 z-20 normal-case">
                    <div className="px-4 py-2 border-b border-slate-200 dark:border-slate-700">
                      <p className="font-semibold text-sm">{user.name}</p>
                      <p className="text-xs text-slate-500 dark:text-slate-400 lowercase">
                        {user.email}
                      </p>
                    </div>
                    <Link
                      href="/profile"
                      onClick={() => {
                        setShowUserMenu(false);
                        setIsOpen(false);
                      }}
                      className="block px-4 py-2 text-sm hover:bg-slate-100 dark:hover:bg-slate-700 transition-colors"
                    >
                      My Profile
                    </Link>
                    <Link
                      href="/bookings"
                      onClick={() => {
                        setShowUserMenu(false);
                        setIsOpen(false);
                      }}
                      className="block px-4 py-2 text-sm hover:bg-slate-100 dark:hover:bg-slate-700 transition-colors"
                    >
                      My Bookings
                    </Link>
                    <button
                      onClick={handleLogout}
                      className="w-full text-left px-4 py-2 text-sm text-red-600 dark:text-red-400 hover:bg-slate-100 dark:hover:bg-slate-700 transition-colors flex items-center gap-2"
                    >
                      <FiLogOut className="text-base" />
                      Logout
                    </button>
                  </div>
                </>
              )}

              <div className="md:hidden flex flex-col items-center gap-4">
                <div className="text-center">
                  <p className="font-semibold normal-case">{user.name}</p>
                  <p className="text-sm text-slate-500 lowercase">
                    {user.email}
                  </p>
                </div>
                <Link
                  href="/profile"
                  onClick={() => setIsOpen(false)}
                  className="text-base"
                >
                  My Profile
                </Link>
                <Link
                  href="/bookings"
                  onClick={() => setIsOpen(false)}
                  className="text-base"
                >
                  My Bookings
                </Link>
                <button
                  onClick={handleLogout}
                  className="text-base text-red-600 dark:text-red-400 font-bold flex items-center gap-2"
                >
                  <FiLogOut className="text-xl" />
                  Logout
                </button>
              </div>
            </div>
          ) : (
            /* Not Authenticated */
            <button
              onClick={handleLogin}
              className="text-foreground font-bold underline hover:opacity-80 transition-opacity"
            >
              Sign In
            </button>
          )}
        </div>
      </div>
    </nav>
  );
};

export default Navbar;
