"use client";

import { useState, useEffect } from "react";
import Stepper, { Step } from "@/components/Stepper";
import { ChevronLeft, ChevronRight, Check, Loader2 } from "lucide-react";
import Navbar from "@/components/Navbar";
import Modal from "@/components/Modal";
import { useBooking } from "@/contexts/BookingContext";
import { useRouter } from "next/dist/client/components/navigation";

interface Timeslot {
  id: number;
  time: string;
  duration: number;
  is_active: boolean;
  available: boolean;
  booked_count: number;
  available_courts: number;
}

interface Court {
  id: number;
  name: string;
  capacity: number;
  description: string;
  is_active: boolean;
  available: boolean;
}

const Book = () => {
  const router = useRouter();
  const { bookingData, setBookingData } = useBooking();
  const API_URL = process.env.NEXT_PUBLIC_API_URL;

  const [selectedDate, setSelectedDate] = useState<string | null>(
    bookingData.selectedDate,
  );
  const [selectedTime, setSelectedTime] = useState<number | null>(
    bookingData.selectedTime,
  );
  const [selectedCourt, setSelectedCourt] = useState<number | null>(
    bookingData.selectedCourt,
  );
  const [currentMonth, setCurrentMonth] = useState(new Date());

  const [detailReservation, setDetailReservation] = useState<{
    id?: number;
    user?: { name: string; email: string };
    court?: Court;
    timeslot?: Timeslot;
  }>({});
  const [openModal, setOpenModal] = useState(false);

  // API Data
  const [availableDates, setAvailableDates] = useState<string[]>([]);
  const [timeslots, setTimeslots] = useState<Timeslot[]>([]);
  const [courts, setCourts] = useState<Court[]>([]);

  // Loading states
  const [loadingDates, setLoadingDates] = useState(true);
  const [loadingTimeslots, setLoadingTimeslots] = useState(false);
  const [loadingCourts, setLoadingCourts] = useState(false);

  const monthNames = [
    "January",
    "February",
    "March",
    "April",
    "May",
    "June",
    "July",
    "August",
    "September",
    "October",
    "November",
    "December",
  ];
  const daysOfWeek = ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"];

  useEffect(() => {
    setBookingData({
      selectedDate,
      selectedTime,
      selectedCourt,
    });
  }, [selectedDate, selectedTime, selectedCourt]);

  // Fetch available dates
  useEffect(() => {
    const fetchDates = async () => {
      try {
        setLoadingDates(true);
        const response = await fetch(`${API_URL}/api/v1/dates`);
        const { data } = await response.json();
        setAvailableDates(data.dates || []);
      } catch (error) {
        console.error("Error fetching dates:", error);
      } finally {
        setLoadingDates(false);
      }
    };

    fetchDates();
  }, []);

  // Fetch timeslots when date is selected
  useEffect(() => {
    if (!selectedDate) return;

    const fetchTimeslots = async () => {
      try {
        setLoadingTimeslots(true);
        const response = await fetch(
          `${API_URL}/api/v1/timeslots?date=${selectedDate}`,
        );
        const { data } = await response.json();
        console.log(data.timeslots);
        setTimeslots(data.timeslots || []);
      } catch (error) {
        console.error("Error fetching timeslots:", error);
      } finally {
        setLoadingTimeslots(false);
      }
    };

    fetchTimeslots();
  }, [selectedDate]);

  // Fetch courts when timeslot is selected
  useEffect(() => {
    if (!selectedDate || !selectedTime) return;

    const fetchCourts = async () => {
      try {
        setLoadingCourts(true);
        const response = await fetch(
          `${API_URL}/api/v1/courts?date=${selectedDate}&timeslot_id=${selectedTime}`,
        );
        const { data } = await response.json();
        setCourts(data.courts || []);
      } catch (error) {
        console.error("Error fetching courts:", error);
      } finally {
        setLoadingCourts(false);
      }
    };

    fetchCourts();
  }, [selectedDate, selectedTime]);

  const handleFinalStepCompleted = async () => {
    // Save to localStorage before navigating
    setBookingData({
      selectedDate,
      selectedTime,
      selectedCourt,
    });

    try {
      if (!selectedDate || !selectedTime || !selectedCourt) return;

      const res = await fetch(`${API_URL}/api/v1/reservations`, {
        method: "POST",
        headers: {
          Authorization: `Bearer ${localStorage.getItem("auth_token")}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          date: selectedDate,
          timeslot_id: selectedTime,
          court_id: selectedCourt,
        }),
      });

      if (!res.ok) {
        console.error("Failed to create reservation");
        return;
      }

      setOpenModal(true);

      const { data } = await res.json();

      setDetailReservation({
        id: data.reservation.ID,
        user: data.reservation.user,
        court: data.reservation.court,
        timeslot: data.reservation.timeslot,
      });

      setOpenModal(true);
    } catch (error) {
      console.error("Error creating reservation:", error);
      return;
    }

    // Navigate to checkout
    // router.push("/checkout");
  };

  const getDaysInMonth = (date: Date) => {
    const year = date.getFullYear();
    const month = date.getMonth();
    const firstDay = new Date(year, month, 1);
    const lastDay = new Date(year, month + 1, 0);
    const daysInMonth = lastDay.getDate();
    const startingDayOfWeek =
      firstDay.getDay() === 0 ? 6 : firstDay.getDay() - 1;
    return { daysInMonth, startingDayOfWeek };
  };

  const prevMonthDays = () => {
    const prevMonth = new Date(
      currentMonth.getFullYear(),
      currentMonth.getMonth(),
      0,
    );
    return prevMonth.getDate();
  };

  const goToPreviousMonth = () => {
    setCurrentMonth(
      new Date(currentMonth.getFullYear(), currentMonth.getMonth() - 1, 1),
    );
  };

  const goToNextMonth = () => {
    setCurrentMonth(
      new Date(currentMonth.getFullYear(), currentMonth.getMonth() + 1, 1),
    );
  };

  const { daysInMonth, startingDayOfWeek } = getDaysInMonth(currentMonth);

  const isDateAvailable = (dateString: string) => {
    return availableDates.includes(dateString);
  };

  const renderCalendarDays = () => {
    const days = [];
    const totalCells = 35;
    const today = new Date();
    today.setHours(0, 0, 0, 0);

    for (let i = 0; i < totalCells; i++) {
      let dayNumber;
      let isCurrentMonth = false;
      let isPrevMonth = false;
      let isNextMonth = false;

      if (i < startingDayOfWeek) {
        dayNumber = prevMonthDays() - startingDayOfWeek + i + 1;
        isPrevMonth = true;
      } else if (i < startingDayOfWeek + daysInMonth) {
        dayNumber = i - startingDayOfWeek + 1;
        isCurrentMonth = true;
      } else {
        dayNumber = i - startingDayOfWeek - daysInMonth + 1;
        isNextMonth = true;
      }

      let cellDate;
      if (isPrevMonth) {
        cellDate = new Date(
          currentMonth.getFullYear(),
          currentMonth.getMonth() - 1,
          dayNumber,
        );
      } else if (isNextMonth) {
        cellDate = new Date(
          currentMonth.getFullYear(),
          currentMonth.getMonth() + 1,
          dayNumber,
        );
      } else {
        cellDate = new Date(
          currentMonth.getFullYear(),
          currentMonth.getMonth(),
          dayNumber,
        );
      }
      cellDate.setHours(0, 0, 0, 0);

      // Format date as YYYY-MM-DD
      const year = cellDate.getFullYear();
      const month = String(cellDate.getMonth() + 1).padStart(2, "0");
      const day = String(cellDate.getDate()).padStart(2, "0");
      const dateString = `${year}-${month}-${day}`;

      const isSelectable = isDateAvailable(dateString);
      const isSelected = selectedDate === dateString;
      const isToday = cellDate.getTime() === today.getTime();

      days.push(
        <button
          key={i}
          onClick={() => {
            if (isSelectable) {
              setSelectedDate(dateString);
              setSelectedTime(null);
              setSelectedCourt(null);
            }
          }}
          disabled={!isSelectable || loadingDates}
          className={`min-h-16 p-2 border border-slate-700 transition-all ${
            isSelectable
              ? "cursor-pointer"
              : "cursor-not-allowed opacity-40 bg-indigo-50/30"
          } ${isSelected ? "bg-indigo-600 " : "hover:bg-indigo-400/30 bg-indigo-50"}
          ${isToday && isCurrentMonth ? "ring-1 ring-indigo-500" : ""}`}
        >
          <div
            className={`text-sm font-bold ${
              isCurrentMonth ? "text-black" : "text-slate-600"
            } ${isSelected ? "text-white" : ""}`}
          >
            {dayNumber}
          </div>
          {isToday && isCurrentMonth && (
            <div className="text-xs text-indigo-400">Today</div>
          )}
        </button>,
      );
    }

    return days;
  };

  const formatDate = (dateString: string | null) => {
    if (!dateString) return "";
    const date = new Date(dateString);
    return date.toLocaleDateString("en-US", {
      weekday: "long",
      year: "numeric",
      month: "long",
      day: "numeric",
    });
  };

  const getSelectedTimeslot = () => {
    return timeslots.find((t) => t.id === selectedTime);
  };

  const getSelectedCourt = () => {
    return courts.find((c) => c.id === selectedCourt);
  };

  return (
    <>
      {openModal && (
        <Modal
          detailReservation={detailReservation}
          handleCloseModal={() => setOpenModal(false)}
        />
      )}

      <Navbar className="text-background bg-white" />
      <section className="h-screen overflow-y-scroll pt-24">
        <div className="mx-auto px-6">
          <div
            className="relative min-h-72 flex items-center justify-center mb-8 bg-cover bg-center rounded-md shadow-xl overflow-hidden"
            style={{ backgroundImage: "url(images/banner-1.jpg)" }}
          >
            <div className="absolute inset-0 bg-background/40"></div>
            <h1 className="relative z-10 text-4xl font-bold text-white -mt-24">
              Pilates Studio Reservation
            </h1>
          </div>

          <Stepper
            initialStep={1}
            onStepChange={(step) => {
              console.log("Current step:", step);
            }}
            onFinalStepCompleted={handleFinalStepCompleted}
            backButtonText="Previous"
            nextButtonText="Next"
            className="-mt-36 mb-16 relative z-10 bg-[#FCF8F8] rounded-xl w-11/12 mx-auto shadow-2xl"
          >
            {/* Step 1: Choose Date */}
            <Step>
              <div className="p-6">
                <h2 className="text-2xl font-bold text-slate-900 dark:text-white mb-6">
                  Select Your Date
                </h2>

                {loadingDates ? (
                  <div className="flex items-center justify-center py-12">
                    <Loader2 className="w-8 h-8 animate-spin text-indigo-600" />
                    <span className="ml-3 text-slate-600">
                      Loading dates...
                    </span>
                  </div>
                ) : (
                  <>
                    <div className="flex items-center justify-between mb-6">
                      <h3 className="text-xl font-semibold text-slate-900 dark:text-white">
                        {monthNames[currentMonth.getMonth()]}{" "}
                        {currentMonth.getFullYear()}
                      </h3>
                      <div className="flex items-center gap-2">
                        <button
                          onClick={goToPreviousMonth}
                          className="p-2 hover:bg-slate-200 dark:hover:bg-slate-800 rounded-lg transition-colors"
                        >
                          <ChevronLeft className="w-5 h-5 text-slate-700 dark:text-slate-300" />
                        </button>
                        <button
                          onClick={goToNextMonth}
                          className="p-2 hover:bg-slate-200 dark:hover:bg-slate-800 rounded-lg transition-colors"
                        >
                          <ChevronRight className="w-5 h-5 text-slate-700 dark:text-slate-300" />
                        </button>
                      </div>
                    </div>

                    {/* Calendar */}
                    <div>
                      <div className="grid grid-cols-7 mb-2">
                        {daysOfWeek.map((day) => (
                          <div
                            key={day}
                            className="text-center text-sm font-medium text-slate-600 dark:text-slate-400 pb-2"
                          >
                            {day}
                          </div>
                        ))}
                      </div>
                      <div className="grid grid-cols-7 gap-1">
                        {renderCalendarDays()}
                      </div>
                    </div>

                    {!selectedDate && (
                      <p className="text-center text-slate-500 dark:text-slate-400 mt-6 text-sm">
                        Please select a date to continue
                      </p>
                    )}
                  </>
                )}
              </div>
            </Step>

            {/* Step 2: Choose Time */}
            <Step>
              <div className="p-6">
                <div className="mb-6">
                  <h2 className="text-2xl font-bold text-slate-900 dark:text-white mb-1">
                    Select Your Time
                  </h2>
                  <p className="text-indigo-600 dark:text-indigo-400">
                    {formatDate(selectedDate)}
                  </p>
                </div>

                {loadingTimeslots ? (
                  <div className="flex items-center justify-center py-12">
                    <Loader2 className="w-8 h-8 animate-spin text-indigo-600" />
                    <span className="ml-3 text-slate-600">
                      Loading time slots...
                    </span>
                  </div>
                ) : (
                  <>
                    <div className="grid grid-cols-3 sm:grid-cols-4 md:grid-cols-6 gap-3">
                      {timeslots
                        .filter((slot) => slot.is_active && slot.available)
                        .map((slot) => (
                          <button
                            key={slot.id}
                            onClick={() => {
                              setSelectedTime(slot.id);
                              setSelectedCourt(null);
                            }}
                            className={`p-4 rounded-lg border-2 transition-all ${
                              selectedTime === slot.id
                                ? "bg-indigo-600 border-indigo-400 text-white"
                                : "bg-slate-100 dark:bg-slate-800 border-slate-300 dark:border-slate-700 text-slate-700 dark:text-slate-300 hover:border-indigo-500 hover:bg-slate-200 dark:hover:bg-slate-700"
                            }`}
                          >
                            <div className="font-medium">{slot.time}</div>
                            <div className="text-xs mt-1 opacity-75">
                              {slot.duration} min
                            </div>
                            {slot.available_courts > 0 && (
                              <div className="text-xs mt-1 opacity-75">
                                {slot.available_courts} courts
                              </div>
                            )}
                          </button>
                        ))}
                    </div>

                    {timeslots.length === 0 && (
                      <p className="text-center text-slate-500 dark:text-slate-400 mt-6 text-sm">
                        No time slots available for this date
                      </p>
                    )}

                    {!selectedTime && timeslots.length > 0 && (
                      <p className="text-center text-slate-500 dark:text-slate-400 mt-6 text-sm">
                        Please select a time to continue
                      </p>
                    )}
                  </>
                )}
              </div>
            </Step>

            {/* Step 3: Choose Court */}
            <Step>
              <div className="p-6">
                <div className="mb-6">
                  <h2 className="text-2xl font-bold text-slate-900 dark:text-white mb-1">
                    Select Your Studio
                  </h2>
                  <p className="text-indigo-600 dark:text-indigo-400">
                    {formatDate(selectedDate)} at {getSelectedTimeslot()?.time}
                  </p>
                </div>

                {loadingCourts ? (
                  <div className="flex items-center justify-center py-12">
                    <Loader2 className="w-8 h-8 animate-spin text-indigo-600" />
                    <span className="ml-3 text-slate-600">
                      Loading studios...
                    </span>
                  </div>
                ) : (
                  <>
                    <div className="grid md:grid-cols-2 gap-4">
                      {courts
                        .filter((court) => court.is_active && court.available)
                        .map((court) => (
                          <button
                            key={court.id}
                            onClick={() => setSelectedCourt(court.id)}
                            className={`p-6 rounded-xl border-2 transition-all text-left ${
                              selectedCourt === court.id
                                ? "bg-indigo-100 dark:bg-indigo-900/50 border-indigo-500"
                                : "bg-slate-100 dark:bg-slate-800/50 border-slate-300 dark:border-slate-700 hover:border-indigo-600"
                            }`}
                          >
                            <div className="flex items-start justify-between mb-3">
                              <h3 className="text-xl font-bold text-slate-900 dark:text-white">
                                {court.name}
                              </h3>
                              {selectedCourt === court.id && (
                                <Check className="w-6 h-6 text-indigo-600 dark:text-indigo-400" />
                              )}
                            </div>
                            <div className="text-sm text-slate-600 dark:text-slate-400 mb-3">
                              Capacity: {court.capacity} people
                            </div>
                            <p className="text-sm text-slate-600 dark:text-slate-400">
                              {court.description}
                            </p>
                          </button>
                        ))}
                    </div>

                    {courts.length === 0 && (
                      <p className="text-center text-slate-500 dark:text-slate-400 mt-6 text-sm">
                        No studios available for this time slot
                      </p>
                    )}

                    {!selectedCourt && courts.length > 0 && (
                      <p className="text-center text-slate-500 dark:text-slate-400 mt-6 text-sm">
                        Please select a studio to complete your reservation
                      </p>
                    )}
                  </>
                )}
              </div>
            </Step>
          </Stepper>
        </div>
      </section>
    </>
  );
};

export default Book;
