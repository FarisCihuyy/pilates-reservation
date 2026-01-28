"use client";

import {
  createContext,
  useContext,
  useState,
  useEffect,
  ReactNode,
} from "react";

interface BookingData {
  selectedDate: string | null;
  selectedTime: number | null;
  selectedCourt: number | null;
}

interface BookingContextType {
  bookingData: BookingData;
  setBookingData: (data: BookingData) => void;
  clearBookingData: () => void;
  saveToSessionStorage: () => void;
  loadFromSessionStorage: () => void;
}

const BookingContext = createContext<BookingContextType | undefined>(undefined);

export function BookingProvider({ children }: { children: ReactNode }) {
  const [bookingData, setBookingDataState] = useState<BookingData>({
    selectedDate: null,
    selectedTime: null,
    selectedCourt: null,
  });

  const STORAGE_KEY = "pilates_booking_data";

  const setBookingData = (data: BookingData) => {
    setBookingDataState(data);
    // Auto save to sessionStorage whenever data changes
    sessionStorage.setItem(STORAGE_KEY, JSON.stringify(data));
  };

  const clearBookingData = () => {
    setBookingDataState({
      selectedDate: null,
      selectedTime: null,
      selectedCourt: null,
    });
    sessionStorage.removeItem(STORAGE_KEY);
  };

  const saveToSessionStorage = () => {
    sessionStorage.setItem(STORAGE_KEY, JSON.stringify(bookingData));
  };

  const loadFromSessionStorage = () => {
    const saved = sessionStorage.getItem(STORAGE_KEY);
    if (saved) {
      try {
        const parsed = JSON.parse(saved);
        setBookingDataState(parsed);
      } catch (error) {
        console.error("Error parsing booking data:", error);
      }
    }
  };

  // Load data on mount
  useEffect(() => {
    loadFromSessionStorage();
  }, []);

  return (
    <BookingContext.Provider
      value={{
        bookingData,
        setBookingData,
        clearBookingData,
        saveToSessionStorage,
        loadFromSessionStorage,
      }}
    >
      {children}
    </BookingContext.Provider>
  );
}

export function useBooking() {
  const context = useContext(BookingContext);
  if (!context) {
    throw new Error("useBooking must be used within BookingProvider");
  }
  return context;
}
