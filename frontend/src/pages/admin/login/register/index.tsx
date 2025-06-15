import RegisterForm from './RegisterForm';

const Register = () => {
  return (
    <section className="flex w-[1280px]">
      <div className="w-[620px] pt-[60px] ">
        <h2 className="mb-10 text-center text-4xl font-bold">
          Регістрація користувача
        </h2>
        <RegisterForm />
      </div>
      <img src="/admin/img-login.png" width="660px" alt="cow" />
    </section>
  );
};

export default Register;
