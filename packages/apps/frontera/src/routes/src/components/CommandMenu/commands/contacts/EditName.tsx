import { useState } from 'react';

import { observer } from 'mobx-react-lite';

import { Edit03 } from '@ui/media/icons/Edit03';
import { useStore } from '@shared/hooks/useStore';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

export const EditName = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;
  const contact = store.contacts.value.get(context.ids?.[0] as string);
  const contactName = contact?.value?.name ?? '';

  const [name, setName] = useState(() => contactName);

  const label = `Contact - ${contact?.value.name}`;

  const handleChangeName = () => {
    if (!context.ids?.[0]) return;

    if (!contact) return;
    contact?.update((o) => {
      o.name = name;

      return o;
    });
    store.ui.commandMenu.setOpen(false);
    store.ui.commandMenu.setType('ContactCommands');
  };

  return (
    <Command>
      <CommandInput
        label={label}
        value={name || ''}
        placeholder='Edit name'
        onValueChange={(value) => setName(value)}
      />
      <Command.List>
        <CommandItem
          leftAccessory={<Edit03 />}
          onSelect={handleChangeName}
        >{`Rename name to "${name}"`}</CommandItem>
      </Command.List>
    </Command>
  );
});
