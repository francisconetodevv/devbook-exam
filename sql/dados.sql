insert into usuarios (nome, nick, email, senha)
values
("usuario1", "usuario_1", "usuario1@gmail.com", "$2a$10$W00f7R4dwqLjrThlU2kF2ely4jGvoheLxtFk9Ov4Qi/orGvlxWncq"),
("usuario2", "usuario_2", "usuario2@gmail.com", "$2a$10$W00f7R4dwqLjrThlU2kF2ely4jGvoheLxtFk9Ov4Qi/orGvlxWncq"),
("usuario3", "usuario_3", "usuario3@gmail.com", "$2a$10$W00f7R4dwqLjrThlU2kF2ely4jGvoheLxtFk9Ov4Qi/orGvlxWncq");

insert into seguidores (usuario_id, seguidor_id)
values
(1, 2),
(3, 1),
(1, 3);

insert into publicacoes (titulo, conteudo, autor_id) 
values 
("Publicacao do usuario 1", "Essa é a publicação do usuario 1! Oba!", 1),
("Publicacao do usuario 2", "Essa é a publicação do usuario 2! Oba!", 2),
("Publicacao do usuario 3", "Essa é a publicação do usuario 3! Oba!", 3);