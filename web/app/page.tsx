import About from "@/components/home/about";
import Hero from "@/components/home/hero";
import Footer from "@/components/layout/footer";
import Header from "@/components/layout/header";

const page = () => {
  return (
    <div>
      <Header />
      <Hero />
      <About />
      <Footer />
    </div>
  );
};

export default page;
