import { cn } from "@/lib/utils";

const H3 = ({
  children,
  className,
}: {
  children: React.ReactNode;
  className: string;
}) => {
  return <h3 className={cn(className)}>{children}</h3>;
};

export default H3;
