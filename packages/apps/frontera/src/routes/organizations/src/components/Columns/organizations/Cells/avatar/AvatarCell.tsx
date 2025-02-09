import { memo, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';

import { useLocalStorage } from 'usehooks-ts';

import { Image } from '@ui/media/Image/Image';
import { Avatar } from '@ui/media/Avatar/Avatar';
import {
  Popover,
  PopoverTrigger,
  PopoverContent,
} from '@ui/overlay/Popover/Popover';

interface AvatarCellProps {
  id: string;
  name: string;
  icon?: string | null;
  logo?: string | null;
  description?: string;
}

export const AvatarCell = memo(
  ({ name, id, icon, logo, description }: AvatarCellProps) => {
    const [isOpen, setIsOpen] = useState(false);
    const navigate = useNavigate();
    const [searchParams] = useSearchParams();
    const preset = searchParams.get('preset');

    const [tabs] = useLocalStorage<{
      [key: string]: string;
    }>(`customeros-player-last-position`, { root: 'organization' });
    const search = searchParams.get('search');
    const [lastSearchForPreset, setLastSearchForPreset] = useLocalStorage<{
      [key: string]: string;
    }>(`customeros-last-search-for-preset`, { root: 'root' });

    const src = icon || logo;
    const fullName = name || 'Unnamed';

    const handleNavigate = () => {
      const lastPositionParams = tabs[id];
      const href = getHref(id, lastPositionParams);

      if (preset) {
        setLastSearchForPreset({
          ...lastSearchForPreset,
          [preset]: search ?? '',
        });
      }
      navigate(href);
    };

    return (
      <div className='items-center ml-[1px]'>
        <Popover open={isOpen} onOpenChange={setIsOpen}>
          <PopoverTrigger>
            <Avatar
              size='xs'
              textSize='xs'
              tabIndex={-1}
              name={fullName}
              src={src || undefined}
              variant='outlineSquare'
              onClick={handleNavigate}
              onMouseEnter={() => setIsOpen(true)}
              onMouseLeave={() => setIsOpen(false)}
              className='text-gray-700 cursor-pointer focus:outline-none'
            />
          </PopoverTrigger>

          <PopoverContent
            className='w-[264px]'
            onCloseAutoFocus={(e) => e.preventDefault()}
          >
            {(logo || icon) && (
              <Image
                src={logo || icon || undefined}
                className='h-[36px] w-fit mb-1'
              />
            )}
            <p className='text-md font-semibold'>{fullName}</p>
            <p className='text-xs'>{description}</p>
          </PopoverContent>
        </Popover>
      </div>
    );
  },
  (prevProps, nextProps) => {
    return (
      prevProps.icon === nextProps.icon && prevProps.logo === nextProps.logo
    );
  },
);

function getHref(id: string, lastPositionParams: string | undefined) {
  return `/organization/${id}?${lastPositionParams || 'tab=about'}`;
}
