import { cn } from "@/lib/utils";

const Loader = ({ className }: { className: string }) => {
  return <div className={cn("",className)}></div>;
};

export default Loader;
