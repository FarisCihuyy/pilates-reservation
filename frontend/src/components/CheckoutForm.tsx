"use client";

import { useEffect, useState } from "react";
import { useBooking } from "@/contexts/BookingContext";

export default function CheckoutForm() {
  const { bookingData, clearBookingData } = useBooking();
  const [formData, setFormData] = useState({
    name: "",
    email: "",
    phone: "",
  });

  useEffect(() => {
    console.log("Booking data in checkout:", bookingData);
  }, [bookingData]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    const reservationData = {
      ...formData,
      date: bookingData.selectedDate,
      timeslot_id: bookingData.selectedTime,
      court_id: bookingData.selectedCourt,
    };

    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/api/v1/reservations`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${localStorage.getItem("auth_token")}`,
          },
          body: JSON.stringify(reservationData),
        },
      );

      if (response.ok) {
        alert("Booking successful!");
        clearBookingData(); // Clear after successful booking
      }
    } catch (error) {
      console.error("Booking error:", error);
    }
  };

  return (
    <div className="max-w-2xl mx-auto p-6">
      <h1 className="text-3xl font-bold mb-6">Complete Your Booking</h1>

      {/* Display booking summary */}
      <div className="bg-slate-100 p-4 rounded-lg mb-6">
        <h2 className="font-semibold mb-2">Booking Summary</h2>
        <p>Date: {bookingData.selectedDate}</p>
        <p>Time: {bookingData.selectedTime}</p>
        <p>Court: {bookingData.selectedCourt}</p>
      </div>

      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block mb-2">Name</label>
          <input
            type="text"
            value={formData.name}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            className="w-full px-4 py-2 border rounded-lg"
            required
          />
        </div>

        <div>
          <label className="block mb-2">Email</label>
          <input
            type="email"
            value={formData.email}
            onChange={(e) =>
              setFormData({ ...formData, email: e.target.value })
            }
            className="w-full px-4 py-2 border rounded-lg"
            required
          />
        </div>

        <div>
          <label className="block mb-2">Phone (Optional)</label>
          <input
            type="tel"
            value={formData.phone}
            onChange={(e) =>
              setFormData({ ...formData, phone: e.target.value })
            }
            className="w-full px-4 py-2 border rounded-lg"
          />
        </div>

        <button
          type="submit"
          className="w-full bg-indigo-600 text-white py-3 rounded-lg font-semibold hover:bg-indigo-700"
        >
          Confirm Booking
        </button>
      </form>
    </div>
  );
}
