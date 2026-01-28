"use client";

import { Check, CreditCard } from "lucide-react";
import { useParams } from "next/navigation";
import Link from "next/link";
import React, { useEffect } from "react";

interface PaymentData {
  court?: {
    name: string;
  };
  reservation?: {
    date: string;
    court?: {
      name: string;
    };
    timeslot?: {
      time: string;
    };
  };
  amount?: number;
}

const PaymentPage = () => {
  const [paymentData, setPaymentData] = React.useState<PaymentData | null>(
    null,
  );
  const [transactionId, setTransactionId] = React.useState<string | null>(null);
  const [isSuccess, setIsSuccess] = React.useState<boolean | null>(null);

  const params = useParams();
  const { reservation_id } = params;

  const id = Number(reservation_id);
  useEffect(() => {
    const fetchPaymentData = async () => {
      try {
        const response = await fetch(
          `${process.env.NEXT_PUBLIC_API_URL}/api/v1/payments/create`,
          {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
              Authorization: `Bearer ${localStorage.getItem("auth_token")}`,
            },
            body: JSON.stringify({ reservation_id: id }),
          },
        );

        if (!response.ok) {
          throw new Error("Failed to fetch payment data");
        }

        const { data } = await response.json();
        const params = new URLSearchParams(data.payment_url.split("?")[1]);
        setTransactionId(params.get("transaction_id"));

        setPaymentData(data.payment);
      } catch (error) {
        const errorMessage =
          error instanceof Error ? error.message : "An error occurred";
        alert(errorMessage);
      }
    };

    fetchPaymentData();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const handleConfirmPayment = async () => {
    const API_URL = process.env.NEXT_PUBLIC_API_URL; // Simulate payment success
    const res = await fetch(`${API_URL}/api/v1/payments/callback`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${localStorage.getItem("auth_token")}`,
      },
      body: JSON.stringify({
        order_id: transactionId,
        transaction_status: "settlement", // SUCCESS
        transaction_id: `MIDTRANS_${transactionId}`,
        status_code: "200",
        gross_amount: paymentData?.amount?.toFixed(2) || "0.00",
      }),
    });

    if (!res.ok) {
      throw new Error("Failed to process payment callback");
    }

    const data = await res.json();
    setIsSuccess(true);
    console.log("Payment success:", data);
  };

  const handleCancelPayment = async () => {
    const API_URL = process.env.NEXT_PUBLIC_API_URL;
    // Simulate payment failed
    await fetch(`${API_URL}/api/v1/payments/callback`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${localStorage.getItem("auth_token")}`,
      },
      body: JSON.stringify({
        order_id: transactionId,
        transaction_status: "deny", // FAILED
        transaction_id: `MIDTRANS_${transactionId}`,
        status_code: "400",
        gross_amount: paymentData?.amount?.toFixed(2) || "0.00",
      }),
    });
    setIsSuccess(false);
    console.log("Payment cancelled");
  };

  const formatDate = (dateString: string) => {
    const options: Intl.DateTimeFormatOptions = {
      year: "numeric",
      month: "long",
      day: "numeric",
    };
    const date = new Date(dateString);
    return date.toLocaleDateString("en-US", options);
  };

  const formatCurrency = (amount: number) => {
    const formatRupiah = new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR",
      minimumFractionDigits: 0,
    }).format(amount);
    return formatRupiah;
  };

  return (
    <>
      {isSuccess && (
        <section className="fixed inset-0 flex items-center justify-center bg-green-100">
          <div className="w-full max-w-2xl bg-white rounded-xl overflow-hidden shadow-md min-h-125">
            <div className="bg-[#022269] text-white p-6">
              <h1 className="mb-1 text-2xl font-bold">Payment Success</h1>
              <p>
                Thank you for your payment. Your booking has been confirmed.
              </p>
            </div>
            <div className="flex flex-col items-center justify-center mx-auto gap-8 p-8 w-2/3 text-background h-full">
              <div className="flex flex-col items-center">
                <div className="block p-2 rounded-full bg-green-800/10">
                  <Check className="text-2xl text-green-800" />
                </div>
                <h2 className="text-xl font-semibold">Payment Success</h2>
              </div>
              <div className="text-sm bg-[#f9fbfc] w-full p-2 rounded-sm shadow">
                <p className="text-sm">
                  Your payment has been successfully processed. You can now
                  proceed to the next step.
                </p>
              </div>
              <div className="flex flex-col items-center">
                <Link
                  href="/reservation"
                  className="inline-block md:text-xl bg-green-600 text-white font-bold px-6 py-2 rounded-md transition-opacity hover:opacity-75"
                >
                  Back to Reservation
                </Link>
              </div>
            </div>
          </div>
        </section>
      )}

      <section className="h-screen flex items-center justify-center bg-[#FDFAF6]">
        <div className="w-full max-w-2xl bg-white rounded-xl overflow-hidden shadow-md min-h-125 ">
          <div className="bg-[#022269] text-white p-6">
            <h1 className="mb-1 text-2xl font-bold">Booking Confirmation</h1>
            <p>
              Please review your booking details before completing the payment
            </p>
          </div>
          <div className="flex flex-col items-center justify-center mx-auto gap-8 p-8 w-2/3 text-background h-full">
            <div className="flex flex-col items-center">
              <div className="block p-2 rounded-full bg-blue-800/10">
                <CreditCard className="text-2xl text-blue-800" />
              </div>
              <h2 className="text-xl font-semibold">Review & Pay</h2>
            </div>
            <div className="text-sm bg-[#f9fbfc] w-full p-2 rounded-sm shadow">
              <h2 className="text-base font-semibold mb-4">Booking Summary</h2>

              <div className="space-y-2">
                <div className="flex justify-between">
                  <h3 className="mb-0.5 font-medium">Service</h3>
                  <p>Pilates Session</p>
                </div>
                <div className="flex justify-between">
                  <h3 className="mb-0.5 font-medium">Location</h3>
                  <p>{paymentData?.reservation?.court?.name}</p>
                </div>
                <div className="flex justify-between">
                  <h3 className="mb-0.5 font-medium">Date</h3>
                  <p>{formatDate(paymentData?.reservation?.date ?? "")}</p>
                </div>
                <div className="flex justify-between">
                  <h3 className="mb-0.5 font-medium">Time</h3>
                  <p>{paymentData?.reservation?.timeslot?.time}</p>
                </div>
              </div>
              <div className="border-t border-t-stone-300/50 pt-4 space-y-2 mt-4">
                <h2 className="text-base font-semibold">Price Details</h2>
                <div className="flex justify-between">
                  <h3 className="mb-0.5 font-medium">Total</h3>
                  <p>{formatCurrency(paymentData?.amount ?? 0)}</p>
                </div>
              </div>
            </div>
            <div className="flex w-full justify-between items-center *:cursor-pointer">
              <button
                onClick={handleCancelPayment}
                className="px-6 py-1 border border-gray-300/50 rounded-sm"
              >
                Cancel
              </button>
              <button
                onClick={handleConfirmPayment}
                className="px-6 py-1 bg-blue-800 text-white rounded-sm"
              >
                Confirm Payment
              </button>
            </div>
          </div>
        </div>
      </section>
    </>
  );
};

export default PaymentPage;
