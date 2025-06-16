const ErrorPage = () => {
  return (
    <div className="relative flex h-screen w-screen items-center justify-center">
      <p className="absolute left-[8%] top-[35%] z-0 text-3xl">
        Щось пійшло не так... Спробуйте пізніше.
      </p>
      <img
        src="/admin/bg.svg"
        alt="cow"
        className="absolute inset-x-0 bottom-0 z-0 w-full"
      />
    </div>
  );
};

export default ErrorPage;
