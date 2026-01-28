type Props = {
  left: React.ReactNode;
  right: React.ReactNode;
};

const AuthWrapper = ({ left, right }: Props) => {
  return (
    <section className="grid grid-cols-2 p-6 min-h-screen gap-6">
      {left}
      {right}
    </section>
  );
};

export default AuthWrapper;
