-- All information in this migration is valid;
-- All users in the migration have:
--  name = testN
--  password = password
--  salt = randomly generated    
INSERT INTO User 
(name, password, salt)
VALUES
(
'test1',
x'a850dfebac5c7b1815ece6f05fc044c6637786a8304099a05b5a75cdadbab56d',
x'59a71f037dcb00eec3622594b0dc5baab25c34ccfa5025872882de622c35050e'
),
(
'test2',
x'f27adc6ce6b514f2f1e4471af21e6b0dd4c66a1c149c71280f604aa312edaba5',
x'634f556b70f7ccb7d66a177f65ea5b3c1aa4aab955ed9b235403b054fe45c8f1'
),
(
'test3',
x'40b98af5e8bee1c765b3e7c557c58493cdc20767834f0f2604353ddcbb01b5d2',
x'94dd46ecbdbc4126598420721de35bcfadbcff740161cc447d8329a904b8f2bd'
);