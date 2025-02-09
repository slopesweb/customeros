import { FC, ReactNode } from 'react';
import { useNavigate } from 'react-router-dom';

import { cn } from '@ui/utils/cn.ts';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@ui/overlay/Popover/Popover.tsx';

interface PaymentDetailsPopoverProps {
  content?: string;
  children: ReactNode;
  withNavigation?: boolean;
}

export const PaymentDetailsPopover: FC<PaymentDetailsPopoverProps> = ({
  withNavigation,
  content,
  children,
}) => {
  const navigate = useNavigate();

  return (
    <Popover>
      <PopoverTrigger>
        <div className='w-full'>{children}</div>
      </PopoverTrigger>
      <PopoverContent
        className={cn(
          content ? 'block' : 'none',
          'w-fit bg-gray-700 text-white rounded-md text-sm border-none',
        )}
      >
        <div className='flex'>
          <p className='text-base mr-2'>{content}</p>

          {withNavigation && (
            <span
              tabIndex={0}
              role='button'
              className={'text-base underline text-white'}
              onClick={() => navigate('/settings?tab=billing')}
            >
              Go to Settings
            </span>
          )}
        </div>
      </PopoverContent>
    </Popover>
  );
};
