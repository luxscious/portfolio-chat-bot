import React, { ReactElement } from "react";

interface AnimatedIconProps {
  href: string;
  size?: number;
  children: ReactElement;
}

const AnimatedIcon: React.FC<AnimatedIconProps> = ({
  href,
  size = 30,
  children,
}) => {
  return (
    <a
      href={href}
      target="_blank"
      rel="noopener noreferrer"
      className="
        inline-flex
        items-center
        justify-center
        transition
        transform
        hover:scale-110
        hover:text-blue-400
      "
    >
      {React.cloneElement(children, { size })}
    </a>
  );
};

export default AnimatedIcon;
