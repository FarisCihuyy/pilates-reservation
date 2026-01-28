import { Rosarivo } from "next/font/google";

const rosarivo = Rosarivo({
  subsets: ["latin"],
  weight: "400",
});

type Props = {
  badge: string;
  title: string;
  subtitle: string;
  image: string;
};

const AuthBanner = ({ badge, title, subtitle, image }: Props) => {
  return (
    <div
      className="text-[#3D0301] flex flex-col justify-between p-8 rounded-2xl bg-cover bg-center"
      style={{ backgroundImage: `url(${image})` }}
    >
      <p className="flex items-center gap-2 uppercase font-light text-lg">
        {badge}
        <span className="inline-block h-px w-24 bg-[#3D0301]" />
      </p>

      <div>
        <h1
          className={`${rosarivo.className} uppercase font-light text-4xl mb-2`}
        >
          {title}
        </h1>
        <p className="opacity-60">{subtitle}</p>
      </div>
    </div>
  );
};

export default AuthBanner;
