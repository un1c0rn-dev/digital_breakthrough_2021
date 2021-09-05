import React from 'react';
import cn from 'classnames';

import s from './Button.module.scss';

interface ButtonInterface {
  id?: string;
  color?: string;
  disabled?: boolean;

  onClick?: (e?: any) => void;

  className?: string;
  children: any;
}

export const Button: React.FC<ButtonInterface> = ({ id, children, disabled = false, onClick, className }) => {
  return (
    <button id={id} disabled={disabled} onClick={onClick} className={cn(s.root, className)}>
      {children}
    </button>
  );
};
