import Link from "next/link";

interface DetailReservation {
  id?: number;
  user?: { name: string; email: string };
  court?: { name: string };
  timeslot?: { time: string; duration: number };
}

const Modal = ({
  detailReservation,
  handleCloseModal,
}: {
  detailReservation: DetailReservation;
  handleCloseModal: () => void;
}) => {
  const { id, user, court, timeslot } = detailReservation;
  console.log(detailReservation);

  return (
    <div className="fixed inset-0 bg-background/40 flex items-center justify-center z-50 text-background">
      <div className="bg-white w-full max-w-96 p-6 rounded-lg shadow-lg">
        <h2 className="text-xl font-bold mb-4 border-b pb-2">
          Detail Reservation
        </h2>

        <div className="space-y-2 text-sm">
          <div className="flex justify-between">
            <h3 className="mb-0.5 font-bold">Name</h3>
            <p>{user?.name}</p>
          </div>
          <div className="flex justify-between">
            <h3 className="mb-0.5 font-bold">Email</h3>
            <p>{user?.email}</p>
          </div>
          <div className="flex justify-between">
            <h3 className="mb-0.5 font-bold">Court Name</h3>
            <p>{court?.name}</p>
          </div>
          <div className="flex justify-between">
            <h3 className="mb-0.5 font-bold">Time</h3>
            <p>{timeslot?.time}</p>
          </div>
          <div className="flex justify-between">
            <h3 className="mb-0.5 font-bold">Duration</h3>
            <p>{timeslot?.duration} min</p>
          </div>
        </div>

        <div className="flex justify-between items-center mt-12">
          <button
            onClick={handleCloseModal}
            className="px-6 py-2 rounded border border-slate-400/50"
          >
            Close
          </button>

          <Link
            href={`/payment/${id}`}
            className="px-6 py-2 rounded bg-green-600 font-bold text-white"
          >
            Continue
          </Link>
        </div>
      </div>
    </div>
  );
};

export default Modal;
