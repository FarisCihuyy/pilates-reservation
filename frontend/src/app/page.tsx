import Navbar from "@/components/Navbar";
import Image from "next/image";
import Link from "next/link";
import { GoArrowUpRight } from "react-icons/go";

export default function Home() {
  return (
    <>
      <Navbar className="text-txt" />
      <section className="h-screen relative w-full grid grid-cols-2 grid-rows-1 bg-background">
        <div className="relative h-full col-span-2 lg:col-span-1">
          <Image
            src="/images/herocp.jpg"
            alt="Hero image"
            fill
            priority
            className="object-cover grayscale object-center"
          />
        </div>
        <div className="absolute inset-0 lg:relative flex items-center justify-center bg-background/80 lg:bg-transparent px-6 md:px-12">
          <div className="space-y-4 text-center md:text-left">
            <h1 className="text-4xl md:text-7xl font-bold uppercase tracking-wide">
              Book Your <span className="text-foreground">Pilates</span> Session
              Easily
            </h1>
            <h3 className=" md:text-2xl">
              Choose your date, time, and studio.
            </h3>
            <div className="flex items-center justify-center md:justify-start gap-8">
              <Link
                href="/reservation"
                className="inline-block md:text-xl font-medium bg-foreground px-6 py-2 rounded-full transition-opacity hover:opacity-75"
              >
                Book a Session
              </Link>
              <Link
                href="/schadule"
                className="group flex items-center gap-2 md:text-xl font-medium text-foreground"
              >
                <span className="relative inline-block after:absolute after:left-0 after:-bottom-1 after:h-px after:w-0 after:bg-foreground after:transition-all group-hover:after:w-full">
                  View Schadule
                </span>
                <GoArrowUpRight className="text-2xl" />
              </Link>
            </div>
          </div>
        </div>
      </section>
    </>
  );
}
