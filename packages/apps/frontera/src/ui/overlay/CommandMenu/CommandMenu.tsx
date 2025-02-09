import { Command, useCommandState } from 'cmdk';

import { cn } from '@ui/utils/cn';
import { Tag, TagLabel } from '@ui/presentation/Tag/Tag';
import { ChevronRight } from '@ui/media/icons/ChevronRight';
import { isUserPlatformMac } from '@utils/getUserPlatform.ts';
import { Command as CommandIcon } from '@ui/media/icons/Command';

interface CommandInputProps
  extends React.InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  value?: string;
  asChild?: boolean;
  placeholder: string;
  children?: React.ReactNode;
  onValueChange?: (value: string) => void;
}

export const CommandInput = ({
  label,
  asChild,
  children,
  placeholder,
  onValueChange,
  ...rest
}: CommandInputProps) => {
  return (
    <div className='relative w-full p-6 pb-2 flex flex-col gap-2 border-b border-b-gray-100'>
      {label && (
        <Tag size='md' variant='subtle' colorScheme='gray'>
          <TagLabel>{label}</TagLabel>
        </Tag>
      )}
      <div className='w-full min-h-10 flex items-center'>
        <Command.Input
          autoFocus
          asChild={asChild}
          children={children}
          placeholder={placeholder}
          onValueChange={onValueChange}
          {...rest}
        />
      </div>
    </div>
  );
};

interface CommandItemProps extends React.HTMLAttributes<HTMLDivElement> {
  disabled?: boolean;
  keywords?: string[];
  onSelect?: () => void;
  children: React.ReactNode;
  leftAccessory?: React.ReactNode;
  rightAccessory?: React.ReactNode;
}

export const CommandItem = ({
  children,
  disabled,
  leftAccessory,
  rightAccessory,
  ...props
}: CommandItemProps) => {
  return (
    <Command.Item disabled={disabled} {...props}>
      {leftAccessory}
      {children}
      <div className='flex gap-1 items-center ml-auto'>{rightAccessory}</div>
    </Command.Item>
  );
};

interface CommandSubItemProps extends Partial<CommandItemProps> {
  leftLabel: string;
  rightLabel: string;
  keywords?: string[];
  icon: React.ReactNode;
  onSelectAction: () => void;
}

export const CommandSubItem = ({
  icon,
  onSelectAction,
  leftLabel,
  rightLabel,
  ...rest
}: CommandSubItemProps) => {
  const search = useCommandState((state) => state.search);

  return (
    <CommandItem
      leftAccessory={icon}
      onSelect={onSelectAction}
      disabled={search.length <= 3}
      className={cn(search.length <= 3 && 'hidden')}
      {...rest}
    >
      <span className='text-gray-500'>{leftLabel}</span>
      <ChevronRight className='mx-1' />
      <span>{rightLabel}</span>
    </CommandItem>
  );
};

export const StaticCommandItem = ({
  children,
  leftAccessory,
  rightAccessory,
  ...props
}: CommandItemProps) => {
  return (
    <div data-cmdk-item {...props}>
      {leftAccessory}
      {children}
      <div className='flex gap-1 items-center ml-auto'>{rightAccessory}</div>
    </div>
  );
};

interface KbdProps extends React.HTMLAttributes<HTMLDivElement> {
  children: React.ReactNode;
}

export const Kbd = ({ children, className, ...props }: KbdProps) => {
  return (
    <kbd
      {...props}
      className={cn(
        'bg-gray-100 text-gray-700 size-5 flex items-center justify-center rounded-[4px] text-xs',
        className,
      )}
    >
      {children}
    </kbd>
  );
};

export const CommandKbd = ({
  className,
}: React.HTMLAttributes<HTMLDivElement>) => {
  if (isUserPlatformMac()) {
    return (
      <Kbd className={className}>
        <CommandIcon className='size-3' />
      </Kbd>
    );
  }

  return (
    <kbd
      className={cn(
        'bg-gray-100 text-gray-700 flex p-1 py-0.5 items-center justify-center rounded-[4px] text-xs',
        className,
      )}
    >
      Ctrl
    </kbd>
  );
};

export { Command, useCommandState };
