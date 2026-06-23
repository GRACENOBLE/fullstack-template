import { cn } from "@/lib/utils";

const H2 = ({
  children,
  className,
}: {
  children: React.ReactNode;
  className?: string;
}) => {
  return <h2 className={cn(className)}>{children}</h2>;
};

export default H2;
