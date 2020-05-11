create table users (
    id bigserial primary key,
    chat_id varchar(300) not null,
    last_video varchar(300),
    start_time int,
    end_time int,
    user_name varchar(300)
);

create unique index users_chat_id_uindex
    on users (chat_id);

create table messages (
    id   serial primary key,
    name varchar,
    text varchar
);

create table settings (
    scale int
)

INSERT INTO public.messages ( name, text) VALUES ('/start', '👋 Привет %s, я делаю гиф из видео!
🎥 Вы можете загрузить любое ваше видео и я сделаю для вас gif. 🔄

❗️ Ваше видео и gif будут храниться в защищённом и безопасном месте , однако настоятельно рекомендуем не присылать видео которые могут каким-то образом скомпрометировать вас! ❗️

❕К сожалению у меня есть некоторые ограничения: ❕
1. Продолжительность gif должно быть не больше 10 секунд.
2. Видео которое вы загружаете должно не превышать 1 минуты.

Вы можете сразу отправить мне ваше любимое видео или ознакомиться с моими командами набрав /help');
INSERT INTO public.messages ( name, text) VALUES ('/cleartime', 'Время сбросилось, введите новое время');
INSERT INTO public.messages ( name, text) VALUES ('download error', 'Не получилось загрузить видео, попробуйте позднее');
INSERT INTO public.messages ( name, text) VALUES ('save video', 'Пожалуйста подождите, сохраняю видео...');
INSERT INTO public.messages ( name, text) VALUES ('successful download', 'Видео успешно загружено, укажите с какой секунды начать делать gif');
INSERT INTO public.messages ( name, text) VALUES ('invalid message', 'Извините, я не понимаю.
Если вы хотите сделать gif, выберите нужный пункт из меню');
INSERT INTO public.messages ( name, text) VALUES ('end second', 'Теперь введите секунду окончания видео');
INSERT INTO public.messages ( name, text) VALUES ('end more start', 'Конец не должен быть меньше начала');
INSERT INTO public.messages ( name, text) VALUES ('video more 10s', 'Продолжительность gif не должно превышать 10 сек');
INSERT INTO public.messages ( name, text) VALUES ('create video', 'Начало и конец получены, начинаем обработку...');
INSERT INTO public.messages ( name, text) VALUES ('start create video', 'Обработка видео завершена...
Началось создание gif...');
INSERT INTO public.messages ( name, text) VALUES ('loading gif', 'Создание gif завершено
Загружаем gif в чат...
');
INSERT INTO public.messages ( name, text) VALUES ('/newgif', 'Отправьте любое видео, выберите время и я сделаю для вас gif');
INSERT INTO public.messages ( name, text) VALUES ('/else', 'Видео останеться то же , просто введите новое время');
INSERT INTO public.messages ( name, text) VALUES ('not video', 'Кажется это не видео.');